<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class ServiceController
{
    /**
     * "/company/",
     * "/company/{id}",
     * "/company/account",
     * "/account",
     * "/account/{id}",
     * "/{id}",
     * "/account/{id}/user",
     * "/tenant/account/blocked",
     */

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function healthCheck(Request $request): JsonResponse
    {
        return new JsonResponse(
            [
                'isHealthy' => true
            ],
            Response::HTTP_OK
        );
    }
}
