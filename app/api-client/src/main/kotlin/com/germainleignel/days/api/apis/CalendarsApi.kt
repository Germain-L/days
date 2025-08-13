package com.germainleignel.days.api.apis

import com.germainleignel.days.api.infrastructure.CollectionFormats.*
import retrofit2.http.*
import retrofit2.Response
import okhttp3.RequestBody
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

import com.germainleignel.days.api.models.HandlersErrorResponse
import com.germainleignel.days.api.models.ServicesCalendarResponse
import com.germainleignel.days.api.models.ServicesCreateCalendarRequest
import com.germainleignel.days.api.models.ServicesUpdateCalendarRequest

interface CalendarsApi {
    /**
     * GET api/calendars
     * Get user calendars
     * Retrieve all calendars for the authenticated user
     * Responses:
     *  - 200: OK
     *  - 401: Unauthorized
     *  - 500: Internal Server Error
     *
     * @return [kotlin.collections.List<ServicesCalendarResponse>]
     */
    @GET("api/calendars")
    suspend fun apiCalendarsGet(): Response<kotlin.collections.List<ServicesCalendarResponse>>

    /**
     * DELETE api/calendars/{id}
     * Delete calendar
     * Delete a calendar and all associated data (user must own the calendar)
     * Responses:
     *  - 204: No Content
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 403: Forbidden
     *  - 404: Not Found
     *  - 500: Internal Server Error
     *
     * @param id Calendar ID
     * @return [Unit]
     */
    @DELETE("api/calendars/{id}")
    suspend fun apiCalendarsIdDelete(@Path("id") id: kotlin.String): Response<Unit>

    /**
     * GET api/calendars/{id}
     * Get calendar by ID
     * Retrieve a specific calendar by ID (user must own the calendar)
     * Responses:
     *  - 200: OK
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 403: Forbidden
     *  - 404: Not Found
     *  - 500: Internal Server Error
     *
     * @param id Calendar ID
     * @return [ServicesCalendarResponse]
     */
    @GET("api/calendars/{id}")
    suspend fun apiCalendarsIdGet(@Path("id") id: kotlin.String): Response<ServicesCalendarResponse>

    /**
     * PUT api/calendars/{id}
     * Update calendar
     * Update calendar name and description (user must own the calendar)
     * Responses:
     *  - 200: OK
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 403: Forbidden
     *  - 404: Not Found
     *  - 409: Conflict
     *  - 500: Internal Server Error
     *
     * @param id Calendar ID
     * @param calendar Calendar update request
     * @return [ServicesCalendarResponse]
     */
    @PUT("api/calendars/{id}")
    suspend fun apiCalendarsIdPut(@Path("id") id: kotlin.String, @Body calendar: ServicesUpdateCalendarRequest): Response<ServicesCalendarResponse>

    /**
     * POST api/calendars
     * Create a new calendar
     * Create a new calendar for the authenticated user
     * Responses:
     *  - 201: Created
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 409: Conflict
     *  - 500: Internal Server Error
     *
     * @param calendar Calendar creation request
     * @return [ServicesCalendarResponse]
     */
    @POST("api/calendars")
    suspend fun apiCalendarsPost(@Body calendar: ServicesCreateCalendarRequest): Response<ServicesCalendarResponse>

}
