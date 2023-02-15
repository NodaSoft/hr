<?php

namespace Db;

class Command
{
    private const TYPE_MAP = [
        'boolean'  => \PDO::PARAM_BOOL,
        'integer'  => \PDO::PARAM_INT,
        'string'   => \PDO::PARAM_STR,
        'resource' => \PDO::PARAM_LOB,
        'NULL'     => \PDO::PARAM_NULL,
    ];
    /**
     * @var string|null
     */
    private ?string $sql;
    /**
     * @var array
     */
    private array $params;
    /**
     * @var \PDOStatement|null
     */
    private ?\PDOStatement $pdoStatement = null;

    public function __construct(?string $sql = null, array $params = [])
    {
        $this->sql = $sql;
        $this->params = $params;
    }

    public function queryAll(): array
    {
        return $this->queryInternal('fetchAll');
    }

    public function queryOne(): ?array
    {
        $result = $this->queryInternal('fetch');
        return  (false === $result) ? null : $result;
    }

    public function execute(): void
    {
        $this->prepare();
        $this->pdoStatement()->execute();
    }

    public function insert()
    {
        $this->execute();
        return Db::i()->pdo()->lastInsertId();
    }

    /**
     * @param array $params
     *
     * @return $this
     */
    public function setParams(array $params): self
    {
        $this->params = $params;
        return $this;
    }

    /**
     * @param array $params
     *
     * @return $this
     */
    public function addParams(array $params): self
    {
        foreach ($params as $param => $value) {
            $this->params[$param] = $value;
        }
        return $this;
    }

    /**
     * @return \PDOStatement
     */
    public function pdoStatement(): \PDOStatement
    {
        return $this->pdoStatement ?? ($this->pdoStatement = Db::i()->pdo()->prepare($this->requireSql()));
    }

    /**
     * @param string $sql
     *
     * @return $this
     */
    public function setSql(string $sql): self
    {
        $this->sql = $sql;
        return $this;
    }

    /**
     * @return string
     */
    public function requireSql(): string
    {
        if (null === $this->sql) {
            throw new \RuntimeException('SQL statement required');
        }

        return $this->sql;
    }

    private function prepare(): void
    {
        $this->bindParams();
    }

    private function queryInternal(string $method)
    {
        $this->prepare();
        $pdoStatement = $this->pdoStatement();
        $result = $pdoStatement->$method(\PDO::FETCH_ASSOC);
        $pdoStatement->closeCursor();

        return $result;
    }

    private function bindParams(): void
    {
        $pdoStatement = $this->pdoStatement();
        foreach ($this->params as $param => $value) {
            $pdoStatement->bindParam($param, $value, $this->paramType($value));
        }
    }

    private function paramType($value)
    {
        $type = \gettype($value);

        return self::TYPE_MAP[$type] ?? \PDO::PARAM_STR;
    }
}