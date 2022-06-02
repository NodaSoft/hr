<?php

namespace Anoshin\Managers;

use PDO;

class DB
{
	public const DSN = 'mysql:dbname=test;host=127.0.0.1';
	public const USER = 'root';
	public const PASSWORD = '';

	private static ?PDO $instance = null;

	public static function getInstance(): PDO
	{
		if (null === self::$instance) {
			self::$instance = new PDO(
				self::DSN,
				self::USER,
				self::PASSWORD
			);
		}

		return self::$instance;
	}
}