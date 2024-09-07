<script lang="ts">
	import { onMount, createEventDispatcher } from "svelte";
	import { baseUrl } from "../config";
	import Spinner from "./Spinner.svelte";
	import type { EmailsListResponse } from "../types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const emailsUrl = baseUrl + "/v1/emails";
	const dispatch = createEventDispatcher();

	let alias = "";
	let email = "";
	let emailOptions: string[] = [];
	let inputElement: HTMLInputElement;
	let validAlias = false;
	let isLoading = false;

	async function handleSubmit() {
		checkValidAlias();

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
				inputElement.focus();
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

	function checkValidAlias() {
		validAlias = inputElement.checkValidity();
		if (!email) {
			validAlias = false;
		}
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
		bind:value={alias}
		on:input={checkValidAlias}
		type="email"
		id="alias"
		name="alias"
		required
	/>

	<select bind:value={email} on:change={checkValidAlias}>
		<option value="">Select e-mail...</option>
		{#each emailOptions as option}
			<option value={option}>{option}</option>
		{/each}
	</select>

	{#if isLoading}
		<Spinner />
	{:else}
		<button type="submit" disabled={validAlias !== true}> Add </button>
	{/if}
</form>

<style></style>
