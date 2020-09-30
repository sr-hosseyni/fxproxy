<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class AppController
{
    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function accounts(Request $request): JsonResponse
    {
        return new JsonResponse(
            [
                'accounts' => [
                    [
                        'company' => 'ABC',
                        'foo' => 'bar 3'
                    ],
                    [
                        'company' => 'XYZ',
                        'foo' => 'bar 1'
                    ],
                    [
                        'company' => 'MNO',
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
    public function account(Request $request, $id): JsonResponse
    {
        return new JsonResponse(
            [
                'account' => [
                    'id' => $id,
                    'foo' => 'bar'
                ]
            ],
            $request->get("hang") ? Response::HTTP_INTERNAL_SERVER_ERROR : Response::HTTP_OK
        );
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function user(Request $request): JsonResponse
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

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function blockTenant(Request $request): JsonResponse
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
