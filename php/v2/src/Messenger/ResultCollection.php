<?php

namespace NodaSoft\Messenger;

/**
 * @implements \Iterator<int, Result>
 */
class ResultCollection implements \Iterator
{
    /** @var Result[] */
    private $collection = [];

    /** @var int */
    private $pointer = 0;

    public function add(Result $result): void
    {
        $this->collection[] = $result;
    }

    /**
     * @return array<int, array<string, mixed>>
     */
    public function toArray(): array
    {
        return array_map(function (Result $result) {
            return $result->toArray();
        }, $this->collection);
    }

    /**
     * @return Result[]
     */
    public function getList(): array
    {
        return $this->collection;
    }

    public function current(): Result
    {
        return $this->collection[$this->pointer];
    }

    public function next(): void
    {
        ++ $this->pointer;
    }

    public function key(): int
    {
        return $this->pointer;
    }

    public function valid(): bool
    {
        return isset($this->collection[$this->pointer]);
    }

    public function rewind(): void
    {
        $this->pointer = 0;
    }
}
