<?php

declare(strict_types=1);

namespace App\Infra;

use App\Domain\User;
use App\Domain\Criteria;
use App\Domain\UserRepository;
use Envms\FluentPDO\Query;

class UserRepositoryPdo implements UserRepository
{
    public function __construct(private Query $query)
    {
    }

    /**
     * @return User[]
     */
    public function getUsersByCriteria(Criteria $criteria, ?int $limit = null): array
    {
        $q = $this->query
            ->from('users');

        if ($criteria->fromAge !== null) {
            $q->where('age >= ?', $criteria->fromAge);
        }
        if ($criteria->names !== null) {
            $q->where('name in', $criteria->names);
        }
        if ($limit !== null) {
            $q->limit($limit);
        }

        return array_map(
            fn($row) => new User(
                id: (int)$row['id'],
                firstName: (string)$row['firstName'],
                lastName: (string)$row['lastName'],
                location: (string)$row['location'],
                age: (int)$row['age'],
                settings: (array)json_decode((string)$row['settings'], true),
            ),
            $q->fetchAll(),
        );
    }

    /**
     * @param User[] $users
     * @return User[]
     */
    public function addUsers(array $users): array
    {
        try {
            $pdo = $this->query->getPdo();
            $pdo->beginTransaction();

            foreach ($users as $user) {
                $id = $this->query->insertInto('users')
                    ->values([
                        'firstName' => $user->firstName,
                        'lastName' => $user->lastName,
                        'age' => $user->age,
                        'location' => $user->location,
                        'settings' => json_encode($user->settings, JSON_THROW_ON_ERROR),
                    ])
                    ->execute();

                $user->id = $id;
            }

            $pdo->commit();

            return $users;
        } catch (\Throwable $th) {
            $pdo->rollBack();
            throw $th;
        }
    }
}
