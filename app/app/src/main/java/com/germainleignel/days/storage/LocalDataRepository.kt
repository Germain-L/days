package com.germainleignel.days.storage

import android.content.Context
import android.content.SharedPreferences
import androidx.compose.ui.graphics.Color
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.getDefaultColors
import com.germainleignel.days.data.model.Calendar
import com.germainleignel.days.data.model.CalendarData
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.withContext
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import java.time.LocalDate

/**
 * Local storage implementation using SharedPreferences and JSON serialization
 */
class LocalDataRepository(context: Context) : DataRepository {

    private val sharedPreferences: SharedPreferences = context.getSharedPreferences(
        PREFS_NAME, Context.MODE_PRIVATE
    )

    private val _settingsFlow = MutableStateFlow(AppSettings())
    private val _calendarDataFlow = MutableStateFlow(CalendarData(emptyList()))

    override val storageType: StorageType = StorageType.LOCAL

    init {
        // Load initial settings and calendar data
        CoroutineScope(Dispatchers.IO).launch {
            _settingsFlow.value = getSettings()
            
            // Migrate existing data if needed
            migrateExistingDataToCalendars()
            
            _calendarDataFlow.value = getCalendarData()
        }
    }

    private suspend fun migrateExistingDataToCalendars() {
        try {
            val existingCalendars = getCalendarData()
            if (existingCalendars.calendars.isEmpty()) {
                // Check if there's existing legacy data to migrate
                val legacyDays = getLegacyColoredDaysForMigration()
                val settings = getSettings()
                
                if (legacyDays.isNotEmpty() || settings.availableColors.isNotEmpty()) {
                    // Create a default calendar with existing settings
                    val defaultCalendar = Calendar.createDefault(
                        name = "My Calendar",
                        colors = settings.availableColors.ifEmpty { getDefaultColors() }
                    )
                    
                    // Save the default calendar
                    val calendarData = CalendarData(
                        calendars = listOf(defaultCalendar),
                        selectedCalendarId = defaultCalendar.id,
                        calendarDays = if (legacyDays.isNotEmpty()) {
                            mapOf(defaultCalendar.id to legacyDays.map { (date, color) ->
                                com.germainleignel.days.data.model.Day(date, color)
                            })
                        } else {
                            emptyMap()
                        }
                    )
                    
                    saveCalendarData(calendarData)
                    
                    // Clear legacy data after migration
                    if (legacyDays.isNotEmpty()) {
                        sharedPreferences.edit().remove(KEY_COLORED_DAYS).apply()
                    }
                }
            }
        } catch (e: Exception) {
            // If migration fails, continue without it
        }
    }

    private fun getLegacyColoredDaysForMigration(): List<Pair<LocalDate, Color>> {
        return try {
            val json = sharedPreferences.getString(KEY_COLORED_DAYS, null) ?: return emptyList()
            val serializedDays = storageJson.decodeFromString<List<SerializableDayColor>>(json)
            serializedDays.map { it.toLocalDateColorPair() }
        } catch (e: Exception) {
            emptyList()
        }
    }

    // Calendar operations
    override suspend fun saveCalendar(calendar: Calendar) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val updatedCalendars = currentData.calendars.filter { it.id != calendar.id } + calendar
            val newData = currentData.copy(calendars = updatedCalendars)
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to save calendar", e)
        }
    }

    override suspend fun deleteCalendar(calendarId: String) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val updatedCalendars = currentData.calendars.filter { it.id != calendarId }
            val updatedDays = currentData.calendarDays.filterKeys { it != calendarId }
            val newSelectedId = if (currentData.selectedCalendarId == calendarId) {
                updatedCalendars.firstOrNull()?.id
            } else {
                currentData.selectedCalendarId
            }
            val newData = currentData.copy(
                calendars = updatedCalendars,
                selectedCalendarId = newSelectedId,
                calendarDays = updatedDays
            )
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to delete calendar", e)
        }
    }

    override suspend fun getCalendars(): List<Calendar> = withContext(Dispatchers.IO) {
        getCalendarData().calendars
    }

    override suspend fun getCalendar(calendarId: String): Calendar? = withContext(Dispatchers.IO) {
        getCalendarData().calendars.find { it.id == calendarId }
    }

    override suspend fun setSelectedCalendar(calendarId: String) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val updatedCalendars = currentData.calendars.map { calendar ->
                calendar.copy(isSelected = calendar.id == calendarId)
            }
            val newData = currentData.copy(
                calendars = updatedCalendars,
                selectedCalendarId = calendarId
            )
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to set selected calendar", e)
        }
    }

    override suspend fun getSelectedCalendar(): Calendar? = withContext(Dispatchers.IO) {
        getCalendarData().getSelectedCalendar()
    }

    override suspend fun getCalendarData(): CalendarData = withContext(Dispatchers.IO) {
        try {
            val json = sharedPreferences.getString(KEY_CALENDAR_DATA, null)
            if (json != null) {
                val serializedData = storageJson.decodeFromString<SerializableCalendarData>(json)
                serializedData.toCalendarData()
            } else {
                // Return empty calendar data if none exists
                CalendarData(emptyList())
            }
        } catch (e: Exception) {
            // If there's an error parsing, return empty data and clear corrupted data
            sharedPreferences.edit().remove(KEY_CALENDAR_DATA).apply()
            CalendarData(emptyList())
        }
    }

    override fun getCalendarDataFlow(): Flow<CalendarData> = _calendarDataFlow.asStateFlow()

    private suspend fun saveCalendarData(calendarData: CalendarData) = withContext(Dispatchers.IO) {
        try {
            val globalSettings = getSettings()
            val serializedData = SerializableCalendarData.fromCalendarData(calendarData, globalSettings)
            val json = storageJson.encodeToString(serializedData)
            sharedPreferences.edit().putString(KEY_CALENDAR_DATA, json).apply()
            _calendarDataFlow.value = calendarData
        } catch (e: Exception) {
            throw StorageException("Failed to save calendar data", e)
        }
    }

    // Calendar-aware day color operations
    override suspend fun saveDayColor(calendarId: String, date: LocalDate, color: Color) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val currentDays = currentData.getDaysForCalendar(calendarId).toMutableList()
            
            // Remove existing day if present, then add new one
            currentDays.removeAll { it.date == date }
            currentDays.add(com.germainleignel.days.data.model.Day(date, color))
            
            val updatedCalendarDays = currentData.calendarDays.toMutableMap()
            updatedCalendarDays[calendarId] = currentDays
            
            val newData = currentData.copy(calendarDays = updatedCalendarDays)
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to save day color", e)
        }
    }

    override suspend fun removeDayColor(calendarId: String, date: LocalDate) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val currentDays = currentData.getDaysForCalendar(calendarId).toMutableList()
            currentDays.removeAll { it.date == date }
            
            val updatedCalendarDays = currentData.calendarDays.toMutableMap()
            updatedCalendarDays[calendarId] = currentDays
            
            val newData = currentData.copy(calendarDays = updatedCalendarDays)
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to remove day color", e)
        }
    }

    override suspend fun getDayColor(calendarId: String, date: LocalDate): Color? = withContext(Dispatchers.IO) {
        getCalendarData().getDaysForCalendar(calendarId).find { it.date == date }?.color
    }

    override suspend fun getAllColoredDays(calendarId: String): Map<LocalDate, Color> = withContext(Dispatchers.IO) {
        getCalendarData().getDaysForCalendar(calendarId).associate { it.date to it.color }
    }

    override suspend fun clearAllDayColors(calendarId: String) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val updatedCalendarDays = currentData.calendarDays.toMutableMap()
            updatedCalendarDays[calendarId] = emptyList()
            
            val newData = currentData.copy(calendarDays = updatedCalendarDays)
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to clear day colors", e)
        }
    }

    override suspend fun saveDayColors(calendarId: String, dayColors: Map<LocalDate, Color>) = withContext(Dispatchers.IO) {
        try {
            val currentData = getCalendarData()
            val days = dayColors.map { (date, color) -> com.germainleignel.days.data.model.Day(date, color) }
            
            val updatedCalendarDays = currentData.calendarDays.toMutableMap()
            updatedCalendarDays[calendarId] = days
            
            val newData = currentData.copy(calendarDays = updatedCalendarDays)
            saveCalendarData(newData)
        } catch (e: Exception) {
            throw StorageException("Failed to save day colors", e)
        }
    }

    // Legacy day color operations (use selected calendar)

    // Legacy day color operations (use selected calendar)
    override suspend fun saveDayColor(date: LocalDate, color: Color) = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            saveDayColor(selectedCalendar.id, date, color)
        } else {
            // Fallback to legacy storage if no calendar is selected
            val currentColors = getLegacyColoredDays().toMutableMap()
            currentColors[date] = color
            saveDayColors(currentColors)
        }
    }

    override suspend fun removeDayColor(date: LocalDate) = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            removeDayColor(selectedCalendar.id, date)
        } else {
            // Fallback to legacy storage if no calendar is selected
            val currentColors = getLegacyColoredDays().toMutableMap()
            currentColors.remove(date)
            saveDayColors(currentColors)
        }
    }

    override suspend fun getDayColor(date: LocalDate): Color? = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            getDayColor(selectedCalendar.id, date)
        } else {
            // Fallback to legacy storage if no calendar is selected
            getLegacyColoredDays()[date]
        }
    }

    override suspend fun getAllColoredDays(): Map<LocalDate, Color> = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            getAllColoredDays(selectedCalendar.id)
        } else {
            // Fallback to legacy storage if no calendar is selected
            getLegacyColoredDays()
        }
    }

    override suspend fun clearAllDayColors() = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            clearAllDayColors(selectedCalendar.id)
        } else {
            // Fallback to legacy storage if no calendar is selected
            sharedPreferences.edit().remove(KEY_COLORED_DAYS).apply()
        }
    }

    // Legacy storage access (for backward compatibility during migration)
    private suspend fun getLegacyColoredDays(): Map<LocalDate, Color> = withContext(Dispatchers.IO) {
        try {
            val json = sharedPreferences.getString(KEY_COLORED_DAYS, null) ?: return@withContext emptyMap()
            val serializedDays = storageJson.decodeFromString<List<SerializableDayColor>>(json)
            serializedDays.associate { it.toLocalDateColorPair() }
        } catch (e: Exception) {
            // If there's an error parsing, return empty map and clear corrupted data
            sharedPreferences.edit().remove(KEY_COLORED_DAYS).apply()
            emptyMap()
        }
    }

    override suspend fun saveDayColors(dayColors: Map<LocalDate, Color>) = withContext(Dispatchers.IO) {
        val selectedCalendar = getSelectedCalendar()
        if (selectedCalendar != null) {
            saveDayColors(selectedCalendar.id, dayColors)
        } else {
            // Fallback to legacy storage if no calendar is selected
            try {
                val serializedDays = dayColors.map { (date, color) ->
                    SerializableDayColor.fromLocalDateColorPair(date, color)
                }
                val json = storageJson.encodeToString(serializedDays)
                sharedPreferences.edit().putString(KEY_COLORED_DAYS, json).apply()
            } catch (e: Exception) {
                throw StorageException("Failed to save day colors", e)
            }
        }
    }

    override suspend fun saveSettings(settings: AppSettings) = withContext(Dispatchers.IO) {
        try {
            val serializedSettings = SerializableAppSettings.fromAppSettings(settings)
            val json = storageJson.encodeToString(serializedSettings)
            sharedPreferences.edit().putString(KEY_SETTINGS, json).apply()
            _settingsFlow.value = settings
        } catch (e: Exception) {
            throw StorageException("Failed to save settings", e)
        }
    }

    override suspend fun getSettings(): AppSettings = withContext(Dispatchers.IO) {
        try {
            val json = sharedPreferences.getString(KEY_SETTINGS, null)
            if (json != null) {
                val serializedSettings = storageJson.decodeFromString<SerializableAppSettings>(json)
                serializedSettings.toAppSettings()
            } else {
                // Return default settings if none exist
                AppSettings()
            }
        } catch (e: Exception) {
            // If there's an error parsing, return default settings and clear corrupted data
            sharedPreferences.edit().remove(KEY_SETTINGS).apply()
            AppSettings()
        }
    }

    override fun getSettingsFlow(): Flow<AppSettings> = _settingsFlow.asStateFlow()

    override suspend fun resetAllData() = withContext(Dispatchers.IO) {
        sharedPreferences.edit().clear().apply()
        _settingsFlow.value = AppSettings()
    }

    /**
     * Export all data as JSON string for backup purposes
     */
    suspend fun exportData(): String = withContext(Dispatchers.IO) {
        val settings = getSettings()
        val coloredDays = getAllColoredDays()
        val exportData = SerializableAppData.fromAppData(settings, coloredDays)
        storageJson.encodeToString(exportData)
    }

    /**
     * Import data from JSON string for restore purposes
     */
    suspend fun importData(jsonData: String): Boolean = withContext(Dispatchers.IO) {
        try {
            val importData = storageJson.decodeFromString<SerializableAppData>(jsonData)
            val (settings, coloredDays) = importData.toAppData()

            saveSettings(settings)
            saveDayColors(coloredDays)
            true
        } catch (e: Exception) {
            false
        }
    }

    companion object {
        private const val PREFS_NAME = "day_tracker_prefs"
        private const val KEY_SETTINGS = "app_settings"
        private const val KEY_COLORED_DAYS = "colored_days"
        private const val KEY_CALENDAR_DATA = "calendar_data"
    }
}

/**
 * Custom exception for storage operations
 */
class StorageException(message: String, cause: Throwable? = null) : Exception(message, cause)
