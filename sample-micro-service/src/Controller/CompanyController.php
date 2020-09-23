<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class CompanyController
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
    public function company(Request $request): JsonResponse
    {
        return new JsonResponse(
            [
                'company' => [
                    'id' => 11,
                    'foo' => 'bar'
                ]
            ],
            Response::HTTP_CREATED
        );
    }
}
