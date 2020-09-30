<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class CompanyController
{
    /**
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
    public function index(Request $request): Response
    {

        return $request->get("hang")
            ? new Response('Internal Service Error!', Response::HTTP_INTERNAL_SERVER_ERROR)
            : new JsonResponse(
                [
                    'companies' => [
                        [
                            'id' => 'acc23849',
                            'foo' => 'bar 3'
                        ],
                        [
                            'id' => 'sd45f768',
                            'foo' => 'bar 1'
                        ],
                        [
                            'id' => 'sj3co3s4',
                            'foo' => 'bar 2'
                        ],
                    ]
                ],
             Response::HTTP_OK
        );
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function show(Request $request, $id): JsonResponse
    {
        return new JsonResponse(
            [
                'company' => [
                    'id' => $id,
                    'foo' => 'bar'
                ]
            ],
            Response::HTTP_OK
        );
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function account(Request $request): JsonResponse
    {
        return new JsonResponse(
            [
                'account' => [
                    'company' => 'sj3co3s4',
                    'foo' => 'bar'
                ]
            ],
            Response::HTTP_OK
        );
    }
}
