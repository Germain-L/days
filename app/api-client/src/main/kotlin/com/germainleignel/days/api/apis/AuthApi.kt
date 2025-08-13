package com.germainleignel.days.api.apis

import com.germainleignel.days.api.infrastructure.CollectionFormats.*
import retrofit2.http.*
import retrofit2.Response
import okhttp3.RequestBody
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

import com.germainleignel.days.api.models.HandlersErrorResponse
import com.germainleignel.days.api.models.ServicesLoginRequest
import com.germainleignel.days.api.models.ServicesLoginResponse

interface AuthApi {
    /**
     * POST api/auth/login
     * User login
     * Authenticate user and return JWT token
     * Responses:
     *  - 200: OK
     *  - 400: Bad Request
     *  - 401: Unauthorized
     *  - 500: Internal Server Error
     *
     * @param credentials Login credentials
     * @return [ServicesLoginResponse]
     */
    @POST("api/auth/login")
    suspend fun apiAuthLoginPost(@Body credentials: ServicesLoginRequest): Response<ServicesLoginResponse>

}
