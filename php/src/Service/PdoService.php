<?php

declare(strict_types=1);

namespace App\Service;

use PDO;

use function _PHPStan_978789531\RingCentral\Psr7\str;

class PdoService
{
    private readonly PDO $connect;
    private array $queries = [];

    public function __construct()
    {
        $this->connect = new PDO(
            getenv('PDO_DSN') ?: '',
            getenv('PDO_USER') ?: null,
            getenv('PDO_PASSWORD') ?: null
        );

        $fileList = glob(__DIR__ . "/../Sql/*.sql");
        $fileList = $fileList ?: [];
        foreach ($fileList as $sqlFile) {
            if (file_exists($sqlFile)) {
                $this->queries[basename($sqlFile, '.sql')] = file_get_contents($sqlFile);
            }
        }
    }

    public function connect(): PDO
    {
        return $this->connect;
    }

    public function exe(string $sql): void
    {
        $this->connect()->prepare($sql)->execute();
    }

    public function fetch(string $query, int $pdoFetch, array $params = []): ?array
    {
        $sth = $this->connect()->prepare($query);
        $sth->execute($params);
        $res = $sth->fetch($pdoFetch);
        return $res ?: null;
    }

    public function fetchAll(
        string $query,
        int $pdoFetch,
        array $params = [],
        int $limit = null,
        int $offset = null
    ): array {
        if ($limit !== null) {
            $query = trim($query, '\n;') . ' LIMIT ' . $limit;
        }

        if ($offset !== null) {
            $query = trim($query, '\n;') . ' OFFSET ' . $offset;
        }

        $sth = $this->connect()->prepare($query);
        $sth->execute($params);
        $res = $sth->fetchAll($pdoFetch);

        return $res ?: [];
    }

    public function getSql(string $name, string $replace = null): string
    {
        $sql = $this->queries[$name];
        if ($replace !== null) {
            $sql = str_replace('_str_replace_', $replace, (string)$sql);
        }
        return $sql;
    }
}
