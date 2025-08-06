package com.germainleignel.days.data

import androidx.compose.ui.graphics.Color

data class AppSettings(
    val selectedColor: Color = getDefaultColors().first().color, // Default to first color
    val isDarkMode: Boolean = false,
    val followSystemTheme: Boolean = true,
    val availableColors: List<ColorWithMeaning> = getDefaultColors()
)

data class ColorWithMeaning(
    val color: Color,
    val meaning: String
)

// Function to get default color palette (3 colors: Bad, Okay, Good)
fun getDefaultColors(): List<ColorWithMeaning> = listOf(
    ColorWithMeaning(Color(0xFFE53E3E), "Bad"), // Red
    ColorWithMeaning(Color(0xFFFFC107), "Okay"), // Yellow/Amber
    ColorWithMeaning(Color(0xFF4CAF50), "Good") // Green
)

// Legacy function for backward compatibility - now uses dynamic colors
fun getDefaultColorMeanings(): Map<Color, String> = getDefaultColors().associate { it.color to it.meaning }

// Legacy color palette - now dynamically generated
val ColorPalette: List<Color>
    get() = getDefaultColors().map { it.color }
