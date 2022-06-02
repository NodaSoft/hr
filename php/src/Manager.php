<?php

namespace Anoshin\Managers;

class Manager
{
    /**
     * Возвращает менеджеров старше заданного возраста.
	 *
     * @param int $age
     * @return array
     */
    public function getOlderThan(int $age): array
    {
		$managers = ManagerRepository::getOlderThan($age);
		$transformed_managers = array();

		if (!empty($managers)) {
			foreach ($managers as $manager) {
				$settings = json_decode($manager['settings'], true);

				$transformed_managers[] = [
					'id' => $manager['id'],
					'name' => $manager['name'],
					'lastName' => $manager['lastName'],
					'from' => $manager['from'],
					'age' => $manager['age'],
					'key' => $settings['key'],
				];
			}
		}

        return $transformed_managers;
    }

	/**
	 * Возвращает менеджеров по списку имен.
	 *
	 * @param $names
	 * @return array
	 */
    public function getByNames($names): array
    {
		if (is_string($names)) {
			$names = explode(',', $names);
		}

		return ManagerRepository::getByNames($names);
    }

	/**
	 * Добавляет менеджеров в базу данных.
	 *
	 * @param array $managers
	 * @return array
	 */
    public function insertAll(array $managers): array
    {
		return ManagerRepository::insertAll($managers);
    }
}