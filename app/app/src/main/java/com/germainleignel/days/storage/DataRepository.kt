package com.germainleignel.days.storage

import androidx.compose.ui.graphics.Color
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.model.Calendar
import com.germainleignel.days.data.model.CalendarData
import kotlinx.coroutines.flow.Flow
import java.time.LocalDate

/**
 * Repository interface that abstracts data storage operations.
 * This allows switching between local storage and remote backend implementations.
 */
interface DataRepository {
    // Calendar operations
    suspend fun saveCalendar(calendar: Calendar)
    suspend fun deleteCalendar(calendarId: String)
    suspend fun getCalendars(): List<Calendar>
    suspend fun getCalendar(calendarId: String): Calendar?
    suspend fun setSelectedCalendar(calendarId: String)
    suspend fun getSelectedCalendar(): Calendar?
    suspend fun getCalendarData(): CalendarData
    fun getCalendarDataFlow(): Flow<CalendarData>

    // Day color operations (calendar-aware)
    suspend fun saveDayColor(calendarId: String, date: LocalDate, color: Color)
    suspend fun removeDayColor(calendarId: String, date: LocalDate)
    suspend fun getDayColor(calendarId: String, date: LocalDate): Color?
    suspend fun getAllColoredDays(calendarId: String): Map<LocalDate, Color>
    suspend fun clearAllDayColors(calendarId: String)

    // Legacy day color operations (for current calendar)
    suspend fun saveDayColor(date: LocalDate, color: Color)
    suspend fun removeDayColor(date: LocalDate)
    suspend fun getDayColor(date: LocalDate): Color?
    suspend fun getAllColoredDays(): Map<LocalDate, Color>
    suspend fun clearAllDayColors()

    // Settings operations (kept for global app settings)
    suspend fun saveSettings(settings: AppSettings)
    suspend fun getSettings(): AppSettings
    fun getSettingsFlow(): Flow<AppSettings>

    // Batch operations for efficiency
    suspend fun saveDayColors(dayColors: Map<LocalDate, Color>)
    suspend fun saveDayColors(calendarId: String, dayColors: Map<LocalDate, Color>)

    // Data management
    suspend fun resetAllData()

    // Storage type identification
    val storageType: StorageType
}

enum class StorageType {
    LOCAL,
    REMOTE
}
