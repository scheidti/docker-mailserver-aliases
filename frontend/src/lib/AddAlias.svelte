<script lang="ts">
	import { onMount, createEventDispatcher } from "svelte";
	import { baseUrl } from "../config";
	import Spinner from "./Spinner.svelte";
	import type { AliasResponse, EmailsListResponse } from "../types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const emailsUrl = baseUrl + "/v1/emails";
	const dispatch = createEventDispatcher();

	let alias = "";
	let domain = "";
	let email = "";
	let emailOptions: string[] = [];
	let inputElement: HTMLInputElement;
	let isLoading = false;
	export let aliases: AliasResponse[] = [];

	$: domainOptions = emailOptions.map((email) => email.split("@")[1]);
	$: aliasAndDomain = alias + "@" + domain;
	$: validAlias =
		alias.length > 0 &&
		domain.length > 0 &&
		email.length > 0 &&
		!checkIfAliasExistsAlready(aliasAndDomain) &&
		(inputElement?.checkValidity() ?? false);

	async function handleSubmit() {
		if (!validAlias) {
			return;
		}

		isLoading = true;

		try {
			const response = await fetch(aliasesUrl, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ alias, email }),
			});

			if (response.status === 201) {
				alias = "";
				email = "";
				domain = "";
				dispatch("added", { alias, email });
			} else {
				console.error(response);
			}
		} catch (error) {
			console.error(error);
		}

		isLoading = false;
		// TODO: Show success or error message
	}

	function checkIfAliasExistsAlready(alias: string) {
		return aliases.some((a) => a.alias === alias);
	}

	async function getEmails() {
		isLoading = true;
		try {
			const response = await fetch(emailsUrl);
			const data: EmailsListResponse = await response.json();
			emailOptions = data.emails;
		} catch {}
		isLoading = false;
	}

	onMount(async () => {
		getEmails();
	});
</script>

<form on:submit|preventDefault={handleSubmit}>
	<input
		bind:this={inputElement}
		bind:value={aliasAndDomain}
		type="email"
		id="aliasAndDomain"
		name="aliasAndDomain"
		class="hidden"
		required
	/>

	<div>
		<label for="alias" class="sr-only">Alias</label>
		<input bind:value={alias} type="text" id="alias" name="alias" />
	</div>

	<div>
		<label for="domain" class="sr-only">Domain</label>
		<select id="domain" name="domain" bind:value={domain}>
			<option value="">Select domain...</option>
			{#each domainOptions as option}
				<option value={option}>{option}</option>
			{/each}
		</select>
	</div>

	<div>
		<label for="email" class="sr-only">Email</label>
		<select id="email" name="email" bind:value={email}>
			<option value="">Select e-mail...</option>
			{#each emailOptions as option}
				<option value={option}>{option}</option>
			{/each}
		</select>
	</div>

	{#if isLoading}
		<Spinner />
	{:else}
		<div>
			<button type="submit" disabled={validAlias !== true}> Add </button>
		</div>
	{/if}
</form>

<style></style>
