package com.germainleignel.days.storage

import androidx.compose.ui.graphics.Color
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
import kotlinx.coroutines.flow.Flow
import java.time.LocalDate

/**
 * Repository interface that abstracts data storage operations.
 * This allows switching between local storage and remote backend implementations.
 */
interface DataRepository {
    // Day color operations
    suspend fun saveDayColor(date: LocalDate, color: Color)
    suspend fun removeDayColor(date: LocalDate)
    suspend fun getDayColor(date: LocalDate): Color?
    suspend fun getAllColoredDays(): Map<LocalDate, Color>
    suspend fun clearAllDayColors()

    // Settings operations
    suspend fun saveSettings(settings: AppSettings)
    suspend fun getSettings(): AppSettings
    fun getSettingsFlow(): Flow<AppSettings>

    // Batch operations for efficiency
    suspend fun saveDayColors(dayColors: Map<LocalDate, Color>)

    // Data management
    suspend fun resetAllData()

    // Storage type identification
    val storageType: StorageType
}

enum class StorageType {
    LOCAL,
    REMOTE
}
