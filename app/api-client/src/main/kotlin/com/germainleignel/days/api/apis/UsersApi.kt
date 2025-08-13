package com.germainleignel.days.api.apis

import com.germainleignel.days.api.infrastructure.CollectionFormats.*
import retrofit2.http.*
import retrofit2.Response
import okhttp3.RequestBody
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

import com.germainleignel.days.api.models.HandlersErrorResponse
import com.germainleignel.days.api.models.ServicesCreateUserRequest
import com.germainleignel.days.api.models.ServicesUserResponse

interface UsersApi {
    /**
     * GET api/users/{id}
     * Get user by ID
     * Retrieve user information by user ID (self only)
     * Responses:
     *  - 200: OK
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 403: Forbidden
     *  - 404: Not Found
     *  - 500: Internal Server Error
     *
     * @param id User ID
     * @return [ServicesUserResponse]
     */
    @GET("api/users/{id}")
    suspend fun apiUsersIdGet(@Path("id") id: kotlin.String): Response<ServicesUserResponse>

    /**
     * POST api/users
     * Create a new user
     * Create a new user account with email and password
     * Responses:
     *  - 201: Created
     *  - 400: Bad Request
     *  - 409: Conflict
     *  - 500: Internal Server Error
     *
     * @param user User creation request
     * @return [ServicesUserResponse]
     */
    @POST("api/users")
    suspend fun apiUsersPost(@Body user: ServicesCreateUserRequest): Response<ServicesUserResponse>

}
