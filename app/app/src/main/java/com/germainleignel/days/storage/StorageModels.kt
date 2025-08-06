package com.germainleignel.days.storage

import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.model.Calendar
import com.germainleignel.days.data.model.CalendarData
import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import java.time.LocalDate

/**
 * Serializable models for storage operations
 */

@Serializable
data class SerializableColor(
    val argb: Long
) {
    fun toColor(): Color = Color(argb.toULong())

    companion object {
        fun fromColor(color: Color): SerializableColor = SerializableColor(color.value.toLong())
    }
}

@Serializable
data class SerializableColorWithMeaning(
    val color: SerializableColor,
    val meaning: String
) {
    fun toColorWithMeaning(): ColorWithMeaning = ColorWithMeaning(color.toColor(), meaning)

    companion object {
        fun fromColorWithMeaning(colorWithMeaning: ColorWithMeaning): SerializableColorWithMeaning =
            SerializableColorWithMeaning(
                SerializableColor.fromColor(colorWithMeaning.color),
                colorWithMeaning.meaning
            )
    }
}

@Serializable
data class SerializableAppColor(
    val color: SerializableColor,
    val meaning: String
) {
    fun toColorWithMeaning(): ColorWithMeaning = ColorWithMeaning(color.toColor(), meaning)

    companion object {
        fun fromColorWithMeaning(colorWithMeaning: ColorWithMeaning): SerializableAppColor =
            SerializableAppColor(
                color = SerializableColor.fromColor(colorWithMeaning.color),
                meaning = colorWithMeaning.meaning
            )
    }
}

@Serializable
data class SerializableCalendar(
    val id: String,
    val name: String,
    val colorScheme: List<SerializableAppColor>,
    val isSelected: Boolean,
    val createdAt: Long
) {
    fun toCalendar(): Calendar = Calendar(
        id = id,
        name = name,
        colorScheme = colorScheme.map { it.toColorWithMeaning() },
        isSelected = isSelected,
        createdAt = createdAt
    )

    companion object {
        fun fromCalendar(calendar: Calendar): SerializableCalendar = SerializableCalendar(
            id = calendar.id,
            name = calendar.name,
            colorScheme = calendar.colorScheme.map { SerializableAppColor.fromColorWithMeaning(it) },
            isSelected = calendar.isSelected,
            createdAt = calendar.createdAt
        )
    }
}

@Serializable
data class SerializableAppSettings(
    val selectedColor: SerializableColor,
    val isDarkMode: Boolean,
    val followSystemTheme: Boolean,
    val availableColors: List<SerializableColorWithMeaning>,
    val hasSeenOnboarding: Boolean = false
) {
    fun toAppSettings(): AppSettings = AppSettings(
        selectedColor = selectedColor.toColor(),
        isDarkMode = isDarkMode,
        followSystemTheme = followSystemTheme,
        availableColors = availableColors.map { it.toColorWithMeaning() },
        hasSeenOnboarding = hasSeenOnboarding
    )

    companion object {
        fun fromAppSettings(settings: AppSettings): SerializableAppSettings = SerializableAppSettings(
            selectedColor = SerializableColor.fromColor(settings.selectedColor),
            isDarkMode = settings.isDarkMode,
            followSystemTheme = settings.followSystemTheme,
            availableColors = settings.availableColors.map { SerializableColorWithMeaning.fromColorWithMeaning(it) },
            hasSeenOnboarding = settings.hasSeenOnboarding
        )
    }
}

@Serializable
data class SerializableDayColor(
    val dateString: String, // ISO format (yyyy-MM-dd)
    val color: SerializableColor
) {
    fun toLocalDateColorPair(): Pair<LocalDate, Color> =
        Pair(LocalDate.parse(dateString), color.toColor())

    companion object {
        fun fromLocalDateColorPair(date: LocalDate, color: Color): SerializableDayColor =
            SerializableDayColor(date.toString(), SerializableColor.fromColor(color))
    }
}

@Serializable
data class SerializableAppData(
    val settings: SerializableAppSettings,
    val coloredDays: List<SerializableDayColor>
) {
    companion object {
        fun fromAppData(settings: AppSettings, coloredDays: Map<LocalDate, Color>): SerializableAppData =
            SerializableAppData(
                settings = SerializableAppSettings.fromAppSettings(settings),
                coloredDays = coloredDays.map { (date, color) ->
                    SerializableDayColor.fromLocalDateColorPair(date, color)
                }
            )
    }

    fun toAppData(): Pair<AppSettings, Map<LocalDate, Color>> {
        val appSettings = settings.toAppSettings()
        val dayColors = coloredDays.associate { it.toLocalDateColorPair() }
        return Pair(appSettings, dayColors)
    }
}

@Serializable
data class SerializableCalendarData(
    val calendars: List<SerializableCalendar>,
    val selectedCalendarId: String? = null,
    val calendarDays: Map<String, List<SerializableDayColor>> = emptyMap(),
    val globalSettings: SerializableAppSettings? = null
) {
    fun toCalendarData(): CalendarData {
        val calendarList = calendars.map { it.toCalendar() }
        val dayColors = calendarDays.mapValues { (_, days) ->
            days.map { com.germainleignel.days.data.model.Day(it.toLocalDateColorPair().first, it.toLocalDateColorPair().second) }
        }
        return CalendarData(
            calendars = calendarList,
            selectedCalendarId = selectedCalendarId,
            calendarDays = dayColors
        )
    }

    companion object {
        fun fromCalendarData(calendarData: CalendarData, globalSettings: AppSettings?): SerializableCalendarData {
            val serializedCalendars = calendarData.calendars.map { SerializableCalendar.fromCalendar(it) }
            val serializedDays = calendarData.calendarDays.mapValues { (_, days) ->
                days.map { SerializableDayColor.fromLocalDateColorPair(it.date, it.color) }
            }
            return SerializableCalendarData(
                calendars = serializedCalendars,
                selectedCalendarId = calendarData.selectedCalendarId,
                calendarDays = serializedDays,
                globalSettings = globalSettings?.let { SerializableAppSettings.fromAppSettings(it) }
            )
        }
    }
}

/**
 * JSON configuration for serialization
 */
val storageJson = Json {
    prettyPrint = true
    ignoreUnknownKeys = true
}
