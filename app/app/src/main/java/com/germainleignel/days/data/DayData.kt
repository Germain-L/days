package com.germainleignel.days.data

import androidx.compose.ui.graphics.Color
import java.time.LocalDate

data class DayData(
    val date: LocalDate,
    val color: Color? = null
)

data class CalendarMonth(
    val year: Int,
    val month: Int,
    val days: List<DayData>
)
