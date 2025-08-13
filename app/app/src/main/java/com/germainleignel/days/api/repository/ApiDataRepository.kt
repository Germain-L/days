package com.germainleignel.days.api.repository

import androidx.compose.ui.graphics.Color
import com.germainleignel.days.api.apis.AuthApi
import com.germainleignel.days.api.apis.CalendarsApi
import com.germainleignel.days.api.apis.UsersApi
import com.germainleignel.days.api.auth.UserSessionManager
import com.germainleignel.days.api.config.ApiConfig
import com.germainleignel.days.api.models.*
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.model.Calendar
import com.germainleignel.days.data.model.CalendarData
import com.germainleignel.days.storage.DataRepository
import com.germainleignel.days.storage.LocalDataRepository
import com.germainleignel.days.storage.StorageType
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.flowOf
import java.time.LocalDate
import java.util.*

/**
 * Repository implementation that uses the Days API for remote data
 * Falls back to local storage for data not yet available via API
 */
class ApiDataRepository(
    private val sessionManager: UserSessionManager,
    private val localRepository: LocalDataRepository
) : DataRepository {

    private var authApi: AuthApi? = null
    private var calendarsApi: CalendarsApi? = null
    private var usersApi: UsersApi? = null

    private val _calendarData = MutableStateFlow(CalendarData(calendars = emptyList()))
    private val calendarDataFlow = _calendarData.asStateFlow()

    override val storageType = StorageType.REMOTE

    init {
        initializeApis()
    }

    private fun initializeApis() {
        val token = sessionManager.getAuthToken()
        val retrofit = ApiConfig.createRetrofit(authToken = token)
        
        authApi = retrofit.create(AuthApi::class.java)
        calendarsApi = retrofit.create(CalendarsApi::class.java)
        usersApi = retrofit.create(UsersApi::class.java)
    }

    // Authentication methods
    suspend fun login(email: String, password: String): Result<UserSessionManager.User> {
        return try {
            sessionManager.setLoading()
            
            val request = ServicesLoginRequest(email = email, password = password)
            val response = authApi?.apiAuthLoginPost(request)
            
            if (response?.isSuccessful == true) {
                val loginResponse = response.body()!!
                val user = UserSessionManager.User(
                    id = loginResponse.user?.id ?: "",
                    email = loginResponse.user?.email ?: "",
                    createdAt = loginResponse.user?.createdAt ?: ""
                )
                
                sessionManager.saveUserSession(loginResponse.token ?: "", user)
                initializeApis() // Reinitialize with new token
                
                Result.success(user)
            } else {
                val error = "Login failed: ${response?.code()}"
                sessionManager.setError(error)
                Result.failure(Exception(error))
            }
        } catch (e: Exception) {
            val error = "Network error: ${e.message}"
            sessionManager.setError(error)
            Result.failure(e)
        }
    }

    suspend fun createUser(email: String, password: String): Result<UserSessionManager.User> {
        return try {
            sessionManager.setLoading()
            
            val request = ServicesCreateUserRequest(email = email, password = password)
            // Use retrofit without auth token for user registration
            val retrofit = ApiConfig.createRetrofit()
            val usersApiNoAuth = retrofit.create(UsersApi::class.java)
            
            val response = usersApiNoAuth.apiUsersPost(request)
            
            if (response.isSuccessful) {
                val userResponse = response.body()!!
                val user = UserSessionManager.User(
                    id = userResponse.id ?: "",
                    email = userResponse.email ?: "",
                    createdAt = userResponse.createdAt ?: ""
                )
                
                Result.success(user)
            } else {
                val error = "Registration failed: ${response.code()}"
                sessionManager.setError(error)
                Result.failure(Exception(error))
            }
        } catch (e: Exception) {
            val error = "Network error: ${e.message}"
            sessionManager.setError(error)
            Result.failure(e)
        }
    }

    fun logout() {
        sessionManager.clearSession()
        initializeApis()
    }

    // Calendar operations - using API
    override suspend fun saveCalendar(calendar: Calendar) {
        try {
            val request = ServicesCreateCalendarRequest(
                name = calendar.name,
                description = "Calendar created from mobile app"
            )
            
            val response = calendarsApi?.apiCalendarsPost(request)
            if (response?.isSuccessful == true) {
                // Refresh calendar list
                refreshCalendarData()
            } else {
                throw Exception("Failed to save calendar: ${response?.code()}")
            }
        } catch (e: Exception) {
            // Fallback to local storage
            localRepository.saveCalendar(calendar)
        }
    }

    override suspend fun deleteCalendar(calendarId: String) {
        try {
            val response = calendarsApi?.apiCalendarsIdDelete(calendarId)
            if (response?.isSuccessful == true) {
                refreshCalendarData()
            } else {
                throw Exception("Failed to delete calendar: ${response?.code()}")
            }
        } catch (e: Exception) {
            // Fallback to local storage
            localRepository.deleteCalendar(calendarId)
        }
    }

    override suspend fun getCalendars(): List<Calendar> {
        return try {
            val response = calendarsApi?.apiCalendarsGet()
            if (response?.isSuccessful == true) {
                response.body()?.map { apiCalendar ->
                    Calendar(
                        id = apiCalendar.id ?: UUID.randomUUID().toString(),
                        name = apiCalendar.name ?: "Unknown",
                        colorScheme = getDefaultColors(), // For now, use default colors
                        isSelected = false
                    )
                } ?: emptyList()
            } else {
                throw Exception("Failed to get calendars: ${response?.code()}")
            }
        } catch (e: Exception) {
            // Fallback to local storage
            localRepository.getCalendars()
        }
    }

    override suspend fun getCalendar(calendarId: String): Calendar? {
        return try {
            val response = calendarsApi?.apiCalendarsIdGet(calendarId)
            if (response?.isSuccessful == true) {
                val apiCalendar = response.body()!!
                Calendar(
                    id = apiCalendar.id ?: calendarId,
                    name = apiCalendar.name ?: "Unknown",
                    colorScheme = getDefaultColors(),
                    isSelected = false
                )
            } else {
                throw Exception("Failed to get calendar: ${response?.code()}")
            }
        } catch (e: Exception) {
            // Fallback to local storage
            localRepository.getCalendar(calendarId)
        }
    }

    private suspend fun refreshCalendarData() {
        val calendars = getCalendars()
        val selectedCalendar = getSelectedCalendar()
        _calendarData.value = CalendarData(
            calendars = calendars,
            selectedCalendarId = selectedCalendar?.id
        )
    }

    // Day color operations - delegated to local storage until API is available
    override suspend fun saveDayColor(calendarId: String, date: LocalDate, color: Color) {
        localRepository.saveDayColor(calendarId, date, color)
    }

    override suspend fun removeDayColor(calendarId: String, date: LocalDate) {
        localRepository.removeDayColor(calendarId, date)
    }

    override suspend fun getDayColor(calendarId: String, date: LocalDate): Color? {
        return localRepository.getDayColor(calendarId, date)
    }

    override suspend fun getAllColoredDays(calendarId: String): Map<LocalDate, Color> {
        return localRepository.getAllColoredDays(calendarId)
    }

    override suspend fun clearAllDayColors(calendarId: String) {
        localRepository.clearAllDayColors(calendarId)
    }

    // Legacy methods delegated to local storage
    override suspend fun saveDayColor(date: LocalDate, color: Color) {
        localRepository.saveDayColor(date, color)
    }

    override suspend fun removeDayColor(date: LocalDate) {
        localRepository.removeDayColor(date)
    }

    override suspend fun getDayColor(date: LocalDate): Color? {
        return localRepository.getDayColor(date)
    }

    override suspend fun getAllColoredDays(): Map<LocalDate, Color> {
        return localRepository.getAllColoredDays()
    }

    override suspend fun clearAllDayColors() {
        localRepository.clearAllDayColors()
    }

    // Calendar selection - delegated to local storage for now
    override suspend fun setSelectedCalendar(calendarId: String) {
        localRepository.setSelectedCalendar(calendarId)
    }

    override suspend fun getSelectedCalendar(): Calendar? {
        return localRepository.getSelectedCalendar()
    }

    override suspend fun getCalendarData(): CalendarData {
        return localRepository.getCalendarData()
    }

    override fun getCalendarDataFlow(): Flow<CalendarData> {
        return localRepository.getCalendarDataFlow()
    }

    // Settings operations - delegated to local storage
    override suspend fun saveSettings(settings: AppSettings) {
        localRepository.saveSettings(settings)
    }

    override suspend fun getSettings(): AppSettings {
        return localRepository.getSettings()
    }

    override fun getSettingsFlow(): Flow<AppSettings> {
        return localRepository.getSettingsFlow()
    }

    // Batch operations - delegated to local storage
    override suspend fun saveDayColors(dayColors: Map<LocalDate, Color>) {
        localRepository.saveDayColors(dayColors)
    }

    override suspend fun saveDayColors(calendarId: String, dayColors: Map<LocalDate, Color>) {
        localRepository.saveDayColors(calendarId, dayColors)
    }

    // Data management
    override suspend fun resetAllData() {
        localRepository.resetAllData()
        sessionManager.clearSession()
    }

    // Helper methods
    private fun getDefaultColors(): List<ColorWithMeaning> {
        return listOf(
            ColorWithMeaning(Color.Red, "Bad Day"),
            ColorWithMeaning(Color.Yellow, "Okay Day"),
            ColorWithMeaning(Color.Green, "Good Day"),
            ColorWithMeaning(Color.Blue, "Great Day")
        )
    }
}
