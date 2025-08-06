package com.germainleignel.days.data.model

data class CalendarData(
    val calendars: List<Calendar>,
    val selectedCalendarId: String? = null,
    val calendarDays: Map<String, List<Day>> = emptyMap() // calendarId -> List<Day>
) {
    fun getSelectedCalendar(): Calendar? {
        return calendars.find { it.id == selectedCalendarId }
    }
    
    fun getDaysForCalendar(calendarId: String): List<Day> {
        return calendarDays[calendarId] ?: emptyList()
    }
    
    fun getCurrentCalendarDays(): List<Day> {
        return selectedCalendarId?.let { getDaysForCalendar(it) } ?: emptyList()
    }
}
