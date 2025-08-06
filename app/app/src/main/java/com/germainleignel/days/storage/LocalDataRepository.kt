package com.germainleignel.days.storage

import android.content.Context
import android.content.SharedPreferences
import androidx.compose.ui.graphics.Color
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.getDefaultColors
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

    override val storageType: StorageType = StorageType.LOCAL

    init {
        // Load initial settings
        CoroutineScope(Dispatchers.IO).launch {
            _settingsFlow.value = getSettings()
        }
    }

    override suspend fun saveDayColor(date: LocalDate, color: Color) = withContext(Dispatchers.IO) {
        val currentColors = getAllColoredDays().toMutableMap()
        currentColors[date] = color
        saveDayColors(currentColors)
    }

    override suspend fun removeDayColor(date: LocalDate) = withContext(Dispatchers.IO) {
        val currentColors = getAllColoredDays().toMutableMap()
        currentColors.remove(date)
        saveDayColors(currentColors)
    }

    override suspend fun getDayColor(date: LocalDate): Color? = withContext(Dispatchers.IO) {
        getAllColoredDays()[date]
    }

    override suspend fun getAllColoredDays(): Map<LocalDate, Color> = withContext(Dispatchers.IO) {
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
        try {
            val serializedDays = dayColors.map { (date, color) ->
                SerializableDayColor.fromLocalDateColorPair(date, color)
            }
            val json = storageJson.encodeToString(serializedDays)
            sharedPreferences.edit().putString(KEY_COLORED_DAYS, json).apply()
        } catch (e: Exception) {
            // Handle serialization errors gracefully
            throw StorageException("Failed to save day colors", e)
        }
    }

    override suspend fun clearAllDayColors() = withContext(Dispatchers.IO) {
        sharedPreferences.edit().remove(KEY_COLORED_DAYS).apply()
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
    }
}

/**
 * Custom exception for storage operations
 */
class StorageException(message: String, cause: Throwable? = null) : Exception(message, cause)
