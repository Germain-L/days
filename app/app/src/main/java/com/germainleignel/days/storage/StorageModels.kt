package com.germainleignel.days.storage

import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
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
data class SerializableAppSettings(
    val selectedColor: SerializableColor,
    val isDarkMode: Boolean,
    val followSystemTheme: Boolean,
    val availableColors: List<SerializableColorWithMeaning>
) {
    fun toAppSettings(): AppSettings = AppSettings(
        selectedColor = selectedColor.toColor(),
        isDarkMode = isDarkMode,
        followSystemTheme = followSystemTheme,
        availableColors = availableColors.map { it.toColorWithMeaning() }
    )

    companion object {
        fun fromAppSettings(settings: AppSettings): SerializableAppSettings = SerializableAppSettings(
            selectedColor = SerializableColor.fromColor(settings.selectedColor),
            isDarkMode = settings.isDarkMode,
            followSystemTheme = settings.followSystemTheme,
            availableColors = settings.availableColors.map { SerializableColorWithMeaning.fromColorWithMeaning(it) }
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

/**
 * JSON configuration for serialization
 */
val storageJson = Json {
    prettyPrint = true
    ignoreUnknownKeys = true
}
