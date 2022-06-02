<?php

namespace Anoshin\Managers;

use Exception;
use PDO;

class ManagerRepository
{
	public const QUERY_LIMIT = 10;

	public static ManagerRepository $instance;

	/**
	 * DB
	 */
	private PDO $db;

	private function __construct(PDO $db) {
		$this->db = $db;
	}

	/**
	 * Реализация singleton
	 */
	public static function getInstance(): ManagerRepository
	{
		if (null === self::$instance) {
			$db = DB::getInstance();
			self::$instance = new self($db);
		}

		return self::$instance;
	}

	/**
	 * Возвращает список менеджеров старше заданного возраста.
	 *
	 * @param int $age
	 * @return array
	 */
	public static function getOlderThan(int $age): array
	{
		$sql = "
			SELECT `id`, `name`, `lastName`, `from`, `age`, `settings`
			FROM `managers`
			WHERE `age` >= {$age}
			LIMIT " . self::QUERY_LIMIT;

		$query = self::getInstance()->db->prepare($sql);
		$query->execute();

		return $query->fetchAll(PDO::FETCH_ASSOC);
	}

	/**
	 * Возвращает менеджеров по именам
	 */
	public static function getByNames(array $names): array
	{
		$managers = array();
		$escaped_names = array();

		// Escape names
		foreach ($names as $name) {
			$name = preg_replace('/[^\w\s]/ui', '', $name);
			$escaped_names[] = $name;
		}

		$escaped_names = array_map(static function ($name) {
			return "'" . $name . "'";
		}, $escaped_names);

		$escaped_names = implode(', ', $escaped_names);

		$sql = "
			SELECT `id`, `name`, `lastName`, `from`, `age`, `settings`
			FROM `managers`
			WHERE `name` IN ({$escaped_names})";

		$query = self::getInstance()->db->prepare($sql);
		$query->execute();

		return $query->fetchAll(PDO::FETCH_ASSOC);
	}

	/**
	 * Добавляет пользователя в базу данных.
	 *
	 * @param string $name
	 * @param string $lastName
	 * @param int $age
	 * @return string
	 */
	public static function add(string $name, string $lastName, int $age): string
	{
		$db = self::getInstance()->db;
		$sth = $db->prepare("INSERT INTO `managers` (`name`, `lastName`, `age`) VALUES (:name, :lastName, :age)");
		$sth->execute([':name' => $name, ':lastName' => $lastName, ':age' => $age]);

		if ($sth->errorCode() !== '00000') {
			throw new \RuntimeException($sth->errorInfo()[2]);
		}

		return $db->lastInsertId();
	}

	/**
	 * Добавляет пользователей в базу данных.
	 *
	 * @param array $managers
	 * @return array
	 */
	public static function insertAll(array $managers): array
	{
		$ids = [];
		$db = self::getInstance()->db;

		$db->beginTransaction();

		try {
			foreach ($managers as $manager) {
				$ids[] = self::add(
					$manager['name'],
					$manager['lastName'],
					$manager['age']
				);
			}

			$db->commit();
		} catch (\Exception $e) {
			$ids = [];
			$db->rollBack();
			print_r($e->getMessage());
		}

		return $ids;
	}
}