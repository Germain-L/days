package com.germainleignel.days.viewmodel

import android.app.Application
import androidx.compose.runtime.mutableStateMapOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.graphics.Color
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorPalette
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.model.Calendar
import com.germainleignel.days.data.model.CalendarData
import com.germainleignel.days.data.getDefaultColors
import com.germainleignel.days.storage.DataRepository
import com.germainleignel.days.storage.RepositoryFactory
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.time.LocalDate

/**
 * Sealed class for different types of errors
 */
sealed class DayTrackerError {
    object StorageError : DayTrackerError()
    object NetworkError : DayTrackerError()
    data class ValidationError(val message: String) : DayTrackerError()
    object UnknownError : DayTrackerError()
}

class DayTrackerViewModel(application: Application) : AndroidViewModel(application) {
    private val repository: DataRepository = RepositoryFactory.getRepository(application)

    // App settings state
    private val _settings = MutableStateFlow(AppSettings())
    val settings = _settings.asStateFlow()

    // Calendar data state
    private val _calendarData = MutableStateFlow(CalendarData(emptyList()))
    val calendarData = _calendarData.asStateFlow()

    // Map to store colored days (date -> color) - cached from repository
    private val _coloredDays = mutableStateMapOf<LocalDate, Color>()
    val coloredDays: Map<LocalDate, Color> = _coloredDays

    // Current selected date for detailed view
    private val _selectedDate = mutableStateOf<LocalDate?>(null)
    val selectedDate = _selectedDate

    // Loading states
    private val _isLoading = MutableStateFlow(false)
    val isLoading = _isLoading.asStateFlow()

    // Error states
    private val _error = MutableStateFlow<DayTrackerError?>(null)
    val error = _error.asStateFlow()

    init {
        loadData()
    }

    private fun loadData() {
        viewModelScope.launch {
            _isLoading.value = true
            _error.value = null
            try {
                // Load settings
                _settings.value = repository.getSettings()

                // Load calendar data
                _calendarData.value = repository.getCalendarData()

                // Load colored days for current calendar
                val savedDays = repository.getAllColoredDays()
                _coloredDays.clear()
                _coloredDays.putAll(savedDays)
            } catch (e: Exception) {
                // Handle loading errors gracefully
                _settings.value = AppSettings()
                _calendarData.value = CalendarData(emptyList())
                _error.value = DayTrackerError.StorageError
            } finally {
                _isLoading.value = false
            }
        }
    }

    // Calendar management methods
    fun createCalendar(name: String, colors: List<ColorWithMeaning>) {
        viewModelScope.launch {
            try {
                val newCalendar = Calendar.createDefault(name, colors)
                repository.saveCalendar(newCalendar)
                _calendarData.value = repository.getCalendarData()
            } catch (e: Exception) {
                _error.value = DayTrackerError.StorageError
            }
        }
    }

    fun deleteCalendar(calendarId: String) {
        viewModelScope.launch {
            try {
                repository.deleteCalendar(calendarId)
                _calendarData.value = repository.getCalendarData()
                // Reload colored days if the deleted calendar was selected
                if (_calendarData.value.selectedCalendarId != calendarId) {
                    val savedDays = repository.getAllColoredDays()
                    _coloredDays.clear()
                    _coloredDays.putAll(savedDays)
                }
            } catch (e: Exception) {
                _error.value = DayTrackerError.StorageError
            }
        }
    }

    fun selectCalendar(calendarId: String) {
        viewModelScope.launch {
            try {
                repository.setSelectedCalendar(calendarId)
                _calendarData.value = repository.getCalendarData()
                
                // Load colored days for the selected calendar
                val savedDays = repository.getAllColoredDays()
                _coloredDays.clear()
                _coloredDays.putAll(savedDays)
            } catch (e: Exception) {
                _error.value = DayTrackerError.StorageError
            }
        }
    }

    fun updateCalendar(calendar: Calendar) {
        viewModelScope.launch {
            try {
                repository.saveCalendar(calendar)
                _calendarData.value = repository.getCalendarData()
            } catch (e: Exception) {
                _error.value = DayTrackerError.StorageError
            }
        }
    }

    fun getSelectedCalendar(): Calendar? {
        return _calendarData.value.getSelectedCalendar()
    }

    fun getCalendars(): List<Calendar> {
        return _calendarData.value.calendars
    }

    fun getCurrentCalendarColors(): List<ColorWithMeaning> {
        return getSelectedCalendar()?.colorScheme ?: getAvailableColorsWithMeanings()
    }

    /**
     * Clear the current error state
     */
    fun clearError() {
        _error.value = null
    }

    /**
     * Retry loading data after an error
     */
    fun retryLoadData() {
        loadData()
    }

    fun toggleDayColor(date: LocalDate) {
        viewModelScope.launch {
            val currentColor = _coloredDays[date]
            if (currentColor != null) {
                // If day is already colored, remove the color
                removeDayColor(date)
            } else {
                // Color the day with the currently selected color
                setDayColor(date, _settings.value.selectedColor)
            }
        }
    }

    fun setDayColor(date: LocalDate, color: Color) {
        viewModelScope.launch {
            try {
                _error.value = null
                if (color == Color.Transparent) {
                    // Transparent color means remove the color
                    repository.removeDayColor(date)
                    _coloredDays.remove(date)
                } else {
                    repository.saveDayColor(date, color)
                    _coloredDays[date] = color
                }
            } catch (e: Exception) {
                // Handle storage errors gracefully
                _error.value = DayTrackerError.StorageError
            }
        }
    }

    fun removeDayColor(date: LocalDate) {
        viewModelScope.launch {
            try {
                repository.removeDayColor(date)
                _coloredDays.remove(date)
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun clearAllColors() {
        viewModelScope.launch {
            try {
                repository.clearAllDayColors()
                _coloredDays.clear()
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun updateSelectedColor(color: Color) {
        viewModelScope.launch {
            try {
                val newSettings = _settings.value.copy(selectedColor = color)
                repository.saveSettings(newSettings)
                _settings.value = newSettings
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun updateDarkMode(isDarkMode: Boolean) {
        viewModelScope.launch {
            try {
                val newSettings = _settings.value.copy(
                    isDarkMode = isDarkMode,
                    followSystemTheme = false // When manually set, don't follow system
                )
                repository.saveSettings(newSettings)
                _settings.value = newSettings
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun setFollowSystemTheme(follow: Boolean) {
        viewModelScope.launch {
            try {
                val newSettings = _settings.value.copy(followSystemTheme = follow)
                repository.saveSettings(newSettings)
                _settings.value = newSettings
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun updateColorMeaning(color: Color, meaning: String) {
        viewModelScope.launch {
            try {
                val updatedColors = _settings.value.availableColors.map { colorWithMeaning ->
                    if (colorWithMeaning.color == color) {
                        colorWithMeaning.copy(meaning = meaning)
                    } else {
                        colorWithMeaning
                    }
                }
                val newSettings = _settings.value.copy(availableColors = updatedColors)
                repository.saveSettings(newSettings)
                _settings.value = newSettings
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun getColorMeaning(color: Color): String {
        return getCurrentCalendarColors().find { it.color == color }?.meaning ?: "Custom Color"
    }

    fun addNewColor(color: Color, meaning: String) {
        viewModelScope.launch {
            try {
                val updatedColors = _settings.value.availableColors + ColorWithMeaning(color, meaning)
                val newSettings = _settings.value.copy(availableColors = updatedColors)
                repository.saveSettings(newSettings)
                _settings.value = newSettings
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun deleteColor(color: Color) {
        viewModelScope.launch {
            try {
                val updatedColors = _settings.value.availableColors.filter { it.color != color }
                var newSettings = _settings.value.copy(availableColors = updatedColors)

                // If deleted color was the selected color, pick the first available color
                if (_settings.value.selectedColor == color && updatedColors.isNotEmpty()) {
                    newSettings = newSettings.copy(selectedColor = updatedColors.first().color)
                }

                repository.saveSettings(newSettings)
                _settings.value = newSettings

                // Remove this color from any colored days
                val daysToRemove = _coloredDays.filter { it.value == color }.keys.toList()
                daysToRemove.forEach { date ->
                    repository.removeDayColor(date)
                    _coloredDays.remove(date)
                }
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun resetAllData() {
        viewModelScope.launch {
            try {
                repository.resetAllData()
                _coloredDays.clear()
                _settings.value = AppSettings()
            } catch (e: Exception) {
                // Handle storage errors gracefully
            }
        }
    }

    fun getDayColor(date: LocalDate): Color? {
        return _coloredDays[date]
    }

    fun getAvailableColors(): List<Color> {
        return getCurrentCalendarColors().map { it.color }
    }

    fun getAvailableColorsWithMeanings(): List<ColorWithMeaning> {
        return _settings.value.availableColors
    }

    // Legacy methods for global colors (used for settings)
    fun getGlobalColors(): List<Color> {
        return _settings.value.availableColors.map { it.color }
    }

    fun getGlobalColorsWithMeanings(): List<ColorWithMeaning> {
        return _settings.value.availableColors
    }

    /**
     * Get current storage type (useful for showing storage status in UI)
     */
    fun getStorageType() = repository.storageType

    /**
     * Export data for backup (if using local storage)
     */
    fun exportData(onResult: (String?) -> Unit) {
        viewModelScope.launch {
            try {
                if (repository is com.germainleignel.days.storage.LocalDataRepository) {
                    val exportedData = repository.exportData()
                    onResult(exportedData)
                } else {
                    onResult(null)
                }
            } catch (e: Exception) {
                onResult(null)
            }
        }
    }

    /**
     * Import data from backup (if using local storage)
     */
    fun importData(jsonData: String, onResult: (Boolean) -> Unit) {
        viewModelScope.launch {
            try {
                if (repository is com.germainleignel.days.storage.LocalDataRepository) {
                    val success = repository.importData(jsonData)
                    if (success) {
                        loadData() // Reload data after import
                    }
                    onResult(success)
                } else {
                    onResult(false)
                }
            } catch (e: Exception) {
                onResult(false)
            }
        }
    }
}
