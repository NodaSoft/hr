<?php

namespace Israil\V2;

use Israil\V2\Interfaces\ReferencesOperationAbstract;
use Israil\V2\Repositories\ClientRepository;
use Israil\V2\Services\Notifications\EmailService;
use Israil\V2\Services\Notifications\PhoneService;
use Israil\V2\Templates\EmailTemplate;
use Israil\V2\Builders\EmailTemplateBuilder;
use Israil\V2\Exceptions\EmailTemplateKeysNotDefined;
use Israil\V2\Validator\EmailTemplateValidator;
use Israil\V2\Repositories\EmployeeRepository;
use Israil\V2\Exceptions\NotFoundException;
use Israil\V2\Enums\NotificationTypeEnum;
use Israil\V2\DTO\OperationDTO;
use Israil\V2\Validator\OperationRequestValidator;
use Israil\V2\Responses\OperationResponse;
use Israil\V2\Repositories\SellerRepository;
use Israil\V2\Stubs\Contractor;
use Israil\V2\Exceptions\InvalidChangedTypeDifferences;
use NW\WebService\References\Operations\Notification\NotificationEvents;

class ReturnOperation extends ReferencesOperationAbstract
{
	protected SellerRepository $sellerRepository;
	protected ClientRepository $clientRepository;
	protected EmployeeRepository $employeeRepository;
	protected EmailService $emailService;
	protected PhoneService $phoneService;

	public function __construct()
	{
		// типа DI
		$this->sellerRepository = new SellerRepository;
		$this->clientRepository = new ClientRepository;
		$this->employeeRepository = new EmployeeRepository;
		$this->emailService = new EmailService;
		$this->phoneService = new PhoneService;
	}

	public function doOperation(): array
	{
		$request = $this->getDtoFromRequest();
		$response = $this->getNewOperationResponseDto();

		try {
			$validatedRequest = OperationRequestValidator::validate($request);
		} catch (NotFoundException $e) {
			return $response->errorResponse($e->getMessage())->toArray();
		}

		try {
			$client = $this->clientRepository->findOrFail($validatedRequest->clientId);
			$expert = $this->employeeRepository->findOrFail($validatedRequest->expertId); // ???
			$employee = $this->employeeRepository->findOrFail($validatedRequest->creatorId);
			$reseller = $this->sellerRepository->findOrFail($validatedRequest->resellerId);
		} catch (NotFoundException $e) {
			return $response->errorResponse($e->getMessage())->toArray();
		}

		try {
			$validatedTemplateData = EmailTemplateValidator::validate($this->buildNewTemplateData(
				$validatedRequest,
				$employee,
				$expert,
				$client,
			));
		} catch (EmailTemplateKeysNotDefined | InvalidChangedTypeDifferences $e) {
			return $response->errorResponse($e->getMessage())->toArray();
		}

		$employeeEmailList = $this->emailService->makeToEmployeeEmailList($validatedTemplateData, $validatedRequest);
		$clientEmailList = $this->emailService->makeToClientEmailList($validatedTemplateData, $validatedRequest, $client);

		$response = $this->sendEmployeeEmail(
			$employeeEmailList,
			$response,
			$reseller->id, NotificationEvents::CHANGE_RETURN_STATUS
		);

		if ($validatedRequest->isType(NotificationTypeEnum::CHANGE)) {
			$response = $this->sendClientEmailIfResellerAndClient(
				$reseller,
				$client,
				$response,
				$validatedTemplateData,
				$clientEmailList,
			);

			$response = $this->sendClientPhoneIfCanReceive(
				$reseller,
				$client,
				$response,
				$validatedTemplateData,
			);
		}

		return $response->toArray();
	}

	protected function sendClientEmailIfResellerAndClient(
		Contractor $seller,
		Contractor $client,
		OperationResponse $response,
		EmailTemplate $template,
		array $emails,
	): OperationResponse {
		if ($this->sellerRepository->getEmailFrom($seller->id) && $client->email) {
			return $this->sendClientEmail(
				$emails,
				$response,
				$seller->id, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $template->DIFFERENCES['to']
			);
		}

		return $response;
	}

	protected function sendClientPhoneIfCanReceive(
		Contractor $seller,
		Contractor $client,
		OperationResponse $response,
		EmailTemplate $template,
	): OperationResponse {
		if ($client->canReceiveSms()) {
			return $this->sendClientPhone(
				$seller,
				$client,
				$response,
				$template,
			);
		}

		return $response;
	}

	private function sendEmployeeEmail(array $emails, OperationResponse $response, ...$args): OperationResponse
	{
		$this->emailService->sendEmailList($emails, ...$args);

		return $response->setIsEmployeeByEmail();
	}

	private function sendClientEmail(array $emails, OperationResponse $response, ...$args): OperationResponse
	{
		$this->emailService->sendEmailList($emails, ...$args);

		return $response->setIsClientByEmail();
	}

	private function sendClientPhone(Contractor $seller, Contractor $client, OperationResponse $response, EmailTemplate $template): OperationResponse
	{
		$smsResponse = $this->phoneService->sendSmsClient(
			$seller->id,
			$client->id,
			NotificationEvents::CHANGE_RETURN_STATUS,
			$template,
		);

		if (!$smsResponse->success) {
			return $response->setClientBySmsMessage($smsResponse->message);
		}

		return $response->setClientBySmsIsSent();
	}

	private function buildNewTemplateData(
		OperationDTO $dto,
		Contractor $employee,
		Contractor $expert,
		Contractor $client,
	): EmailTemplate {
		return (new EmailTemplateBuilder(new EmailTemplate))
			->setComplaintId($dto->complaintId)
			->setComplaintNumber($dto->complaintNumber)
			->setClientId($dto->clientId)
			->setCreatorId($dto->creatorId)
			->setExpertId($dto->expertId)
			->setConsumptionId($dto->consumptionId)
			->setCreatorName($employee->getFullName())
			->setExpertName($expert->getFullName())
			->setClientName($client->getFullName())
			->setConsumptionNumber($dto->consumptionNumber)
			->setAgreementNumber($dto->agreementNumber)
			->setDifferences($dto->notificationType, $dto->resellerId, $dto->differences)
			->setDate($dto->date)
			->resolve();
	}

	private function getDtoFromRequest(): OperationDTO
	{
//		return OperationDTO::fromArray($this->getRequest('data'));
		return OperationDTO::fromArray([
			'resellerId' => 1,
			'notificationType' => 1,
			'clientId' => 1,
			'creatorId' => 1,
			'expertId' => 1,
			'complaintId' => 1,
			'complaintNumber' => '123-4567-890',
			'consumptionId' => 1,
			'consumptionNumber' => '098-7654-123',
			'agreementNumber' => '098-7654-123',
			'differences' => [],
			'date' => '12-12-2024',
		]);
	}

	private function getNewOperationResponseDto(): OperationResponse
	{
		return new OperationResponse;
	}
}
