<?php

namespace Gateway;

use Db\Command;

class UsersGateway
{
    private array $where = [];

    /**
     * Добавляет пользователя в базу данных.
     *
     * @param string $name
     * @param string $lastName
     * @param int    $age
     *
     * @return string
     */
    public function add(string $name, string $lastName, int $age): string
    {
        return (new Command())
            ->setSql('INSERT INTO `user` (`name`, `last_name`, `age`, `settings`) VALUES (:name, :age, :lastName, :settings)')
            ->setParams([
                ':name'     => $name,
                ':age'      => $age,
                ':lastName' => $lastName,
                ':settings' => '{"keep":"keep"}',
            ])
            ->insert()
        ;
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     *
     * @param int $ageFrom
     *
     * @return $this
     */
    public function byAge(int $ageFrom): self
    {
        $this->where['age>'] = ['>', 'age', [$ageFrom]];

        return $this;
    }

    /**
     * Возвращает пользователя по имени.
     *
     * @param string|string[] $name
     *
     * @return $this
     */
    public function byName($name): self
    {
        $name = (array)$name;

        if ([] === $name) {
            throw new \RuntimeException('Empty Names List');
        }

        $this->where['name='] = [(\count($name) > 1) ? 'IN' : '=', 'name', $name];

        return $this;
    }

    public function all(?int $limit = null): array
    {
        $result = $this->buildCommand($limit)->queryAll();
        \array_walk($result, function (&$value) {
            $value = $this->map($value);
        });
        return $result;
    }

    public function one(): ?array
    {
        $result = $this->buildCommand(1)->queryOne();
        return (null === $result) ? null : $this->map($result);
    }

    private function map($row): array
    {
        try {
            $settings = \json_decode($row['settings'], true, 512, \JSON_THROW_ON_ERROR);
        } catch (\JsonException $jsonException) {
            $settings = [];
        }
        return [
            'id'       => $row['id'],
            'name'     => $row['name'],
            'lastName' => $row['last_name'],
            'from'     => $row['address'],
            'age'      => $row['age'],
            'key'      => $settings['key'] ?? null,
        ];
    }

    private function buildCommand(?int $limit): Command
    {
        $sql = 'SELECT * FROM `user`';

        $conditions = [];
        $params = [];
        $index = 0;
        foreach ($this->where as $item) {
            $conditionParams = [];
            foreach ($item[2] as $value) {
                $paramName = ':p' . ++$index;
                $conditionParams[] = $paramName;
                $params[$paramName] = $value;
            }
            $conditionParams = (1 === \count($conditionParams))
                ? \reset($conditionParams)
                : '(' . \implode(', ', $conditionParams) . ')';
            $conditions[] = "`{$item[1]}` {$item[0]} {$conditionParams}";
        }

        if ([] !== $conditions) {
            $sql .= ' WHERE ' . \implode(' AND ', $conditions);
        }

        if (null !== $limit) {
            $sql .= " LIMIT {$limit}";
        }

        return (new Command())
            ->setSql($sql)
            ->setParams($params)
        ;
    }
}