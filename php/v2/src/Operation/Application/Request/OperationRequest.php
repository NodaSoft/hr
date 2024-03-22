<?php

namespace src\Operation\Application\Request;

use src\Operation\Application\Exceptions\ValidationException;
use src\Operation\Infrastructure\Domain\Enum\NotificationType;

class OperationRequest
{
    private array $data;

    public function __construct(array $request)
    {
        $this->data = $request;
    }

    /**
     * Валидирует данные запроса и возвращает их, если они корректны.
     *
     * @return array Валидированные данные.
     * @throws ValidationException Если данные не прошли валидацию.
     */
    public function validated(): array
    {
        $errors = $this->validate($this->data);

        if (!empty($errors)) {
            throw new ValidationException("Validation failed: " . implode('; ', $errors), 400);
        }

        return $this->data;
    }

    /**
     * Проверяет данные запроса на корректность.
     *
     * @param array $data Данные для валидации.
     * @return array Список ошибок.
     */
    private function validate(array $data): array
    {
        $errors = [];

        if (empty($data['resellerId'])) {
            $errors[] = 'Empty resellerId';
        }

        if (empty($data['clientId'])) {
            $errors[] = 'Empty clientId';
        }

        if (empty($data['notificationType'])
            || !in_array(
                $data['notificationType'],
                [
                    NotificationType::TYPE_NEW,
                    NotificationType::TYPE_CHANGE
                ]
            )
        ) {
            $errors[] = 'Invalid or empty notificationType';
        }

        if (isset($data['creatorId']) && !is_numeric($data['creatorId'])) {
            $errors[] = 'Invalid creatorId';
        }

        if (isset($data['expertId']) && !is_numeric($data['expertId'])) {
            $errors[] = 'Invalid expertId';
        }

        if (isset($data['difference']) && !is_array($data['difference']) ||
            !key_exists('from', $data['difference']) ||
            !key_exists('to', $data['difference'])) {

            $errors[] = 'Invalid differences';
        }

        return $errors;
    }
}