<?php

namespace Israil\V2\Services\Notifications;

use Israil\V2\Dto\OperationDTO;
use Israil\V2\Repositories\SellerRepository;
use Israil\V2\Stubs\MessagesClient;
use Israil\V2\Templates\EmailTemplate;
use Israil\V2\Stubs\Contractor;
use function NW\WebService\References\Operations\Notification\getEmailsByPermit;

class EmailService
{
	public SellerRepository $sellerRepository;

	public function __construct()
	{
		$this->sellerRepository = new SellerRepository;
	}

	public function makeToEmployeeEmailList(EmailTemplate $template, OperationDTO $dto): array
	{
		$emailFrom = $this->sellerRepository->getEmailFrom($dto->resellerId);
		$emailList = $this->getEmployeeEmailsFromConfig($dto->resellerId);

		if (!$emailFrom && !empty($emailList)) {
			return [];
		}

		$serializedEmails = [];
		$templateData = $template->toArray();

		foreach ($emailList as $email) {
			$serializedEmails[] = $this->makeEmailBody($emailFrom, $email, $templateData, $dto->resellerId);
		}

		return $serializedEmails;
	}

	public function makeToClientEmailList(EmailTemplate $template, OperationDTO $dto, Contractor $client): array
	{
		$emailFrom = $this->sellerRepository->getEmailFrom($dto->resellerId);

		if (!$emailFrom && !$client->email) {
			return [];
		}

		return [$this->makeEmailBody($emailFrom, $client->email, $template->toArray(), $dto->resellerId)];
	}

	public function sendEmailList(array $emailList, ...$args): bool
	{
		return MessagesClient::sendMessage($emailList, ...$args);
	}

	protected function getEmployeeEmailsFromConfig(int $resellerId): array
	{
		return getEmailsByPermit($resellerId, 'tsGoodsReturn');
	}

	protected function makeEmailBody(string $from, string $to, array $template, int $resellerId): array
	{
		return [
			'emailFrom' => $from,
			'emailTo' => $to,
			'subject' => __('complaintEmployeeEmailSubject', $template, $resellerId),
			'message' => __('complaintEmployeeEmailBody', $template, $resellerId),
		];
	}
}