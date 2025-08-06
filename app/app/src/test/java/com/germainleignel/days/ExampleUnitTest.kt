package com.germainleignel.days

import androidx.compose.ui.graphics.Color
import com.germainleignel.days.data.AppSettings
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.getDefaultColors
import org.junit.Test
import org.junit.Assert.*
import java.time.LocalDate

/**
 * Unit tests for the Days app core functionality
 */
class DaysUnitTest {

    @Test
    fun appSettings_defaultValues_areCorrect() {
        val defaultSettings = AppSettings()
        
        assertEquals(false, defaultSettings.isDarkMode)
        assertEquals(true, defaultSettings.followSystemTheme)
        assertEquals(3, defaultSettings.availableColors.size)
        assertEquals(getDefaultColors().first().color, defaultSettings.selectedColor)
    }

    @Test
    fun colorWithMeaning_creation_isCorrect() {
        val red = Color(0xFFE53E3E)
        val colorWithMeaning = ColorWithMeaning(red, "Bad")
        
        assertEquals(red, colorWithMeaning.color)
        assertEquals("Bad", colorWithMeaning.meaning)
    }

    @Test
    fun defaultColors_haveCorrectMeanings() {
        val defaultColors = getDefaultColors()
        
        assertEquals(3, defaultColors.size)
        assertEquals("Bad", defaultColors[0].meaning)
        assertEquals("Okay", defaultColors[1].meaning)
        assertEquals("Good", defaultColors[2].meaning)
    }

    @Test
    fun localDate_creation_isCorrect() {
        val date = LocalDate.of(2025, 8, 6)
        
        assertEquals(2025, date.year)
        assertEquals(8, date.monthValue)
        assertEquals(6, date.dayOfMonth)
    }

    @Test
    fun color_argb_conversion_works() {
        val red = Color(0xFFE53E3E)
        val argbValue = red.value
        
        assertTrue(argbValue > 0u)
        assertEquals(Color(argbValue), red)
    }
}