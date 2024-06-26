<?php

namespace Nodasoft\Testapp\Http\Responses;

use Nyholm\Psr7\Factory\Psr17Factory;
use Psr\Http\Message\ResponseInterface;

class JsonResponse
{
    public static function json(array $data, string $message = '', int $status = 200): ResponseInterface
    {
        $psr17Factory = new Psr17Factory();
        $response = $psr17Factory->createResponse($status);
        // Set Content-Type header
        $response->withHeader('Content-Type', 'application/json');
        // Write JSON data to response body
        $response->getBody()->write(json_encode([
            'message' => $message,
            'data' => $data
        ]));
        return $response;
    }
}