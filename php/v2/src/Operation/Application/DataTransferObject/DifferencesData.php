<?php

namespace Src\Operation\Application\DataTransferObject;

class DifferencesData
{
    public int $from;
    public ?int $to;

    public function __construct(int $from, int $to)
    {
        $this->from = $from;
        $this->to = $to;
    }

    /**
     * @return array
     */
    public function toArray(): array
    {
        return [
            'from' => $this->from,
            'to' => $this->to ?? null,
        ];
    }

    /**
     * @param array $data
     * @return self
     */
    public static function fromArray(array $data): self
    {
        return new self($data['from'], $data['to'] ?? null);
    }

}