<script lang="ts">
	import { onMount } from "svelte";
	import { baseUrl } from "../config";
	import { toasts } from "../stores";
	import Spinner from "./Spinner.svelte";
	import type { AliasResponse, EmailsListResponse } from "../types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const emailsUrl = baseUrl + "/v1/emails";

	let alias = $state("");
	let domain = $state("");
	let email = $state("");
	let emailOptions: string[] = $state([]);
	let inputElement: HTMLInputElement | undefined = $state();
	let isLoading = $state(false);
	let includeExistingAliases = $state(false);
	interface Props {
		aliases?: AliasResponse[];
		added?: (data: { alias: string; email: string }) => void;
	}

	let { aliases = [], added }: Props = $props();

	async function handleSubmit(event: Event) {
		event.preventDefault();
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
				body: JSON.stringify({ alias: aliasAndDomain, email }),
			});

			if (response.status === 201) {
				alias = "";
				email = "";
				domain = "";
				added?.({ alias, email });
				toasts.update((toasts) => [
					...toasts,
					{ type: "success", text: "Alias added" },
				]);
			} else {
				toasts.update((toasts) => [
					...toasts,
					{
						type: "error",
						text: `Failed to add alias: ${response.statusText}`,
					},
				]);
			}
		} catch (error) {
			toasts.update((toasts) => [
				...toasts,
				{ type: "error", text: `Failed to add alias: ${error}` },
			]);
		}

		isLoading = false;
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

	let domainOptions = $derived((() => {
		const result = emailOptions.map((email) => email.split("@")[1])
		const aliasDomains = aliases.map((alias) => alias.alias.split("@")[1]);
		const aliasEmailDomains = aliases.map((alias) => alias.email.split("@")[1]);
		result.push(...aliasDomains);
		result.push(...aliasEmailDomains);
		return [...new Set(result)].sort((a, b) => a.localeCompare(b));
	})());

	let emailSelectOptions = $derived((() => {
		const result = [...emailOptions];
		if (includeExistingAliases) {
			result.push(...aliases.map((alias) => alias.alias));
		}
		result.sort((a, b) => a.localeCompare(b));
		return result;
	})());

	$effect(() => {
		if (!emailSelectOptions.includes(email)) {
			email = "";
		}
		if (!domainOptions.includes(domain)) {
			domain = "";
		}
	});

	let aliasAndDomain = $derived(alias + "@" + domain);

	let validAlias =
		$derived(alias.length > 0 &&
		domain.length > 0 &&
		email.length > 0 &&
		!checkIfAliasExistsAlready(aliasAndDomain) &&
		(inputElement?.checkValidity() ?? false));
</script>

<div class="mx-auto flex justify-center items-center">
	<form onsubmit={handleSubmit}>
		<input
			bind:this={inputElement}
			bind:value={aliasAndDomain}
			type="email"
			id="aliasAndDomain"
			name="aliasAndDomain"
			class="hidden"
			required
		/>

		<div class="add-row">
			<p class="text-lg font-bold text-primary">Create new alias</p>
		</div>

		<div class="add-row">
			<div>
				<label for="alias" class="sr-only">Alias</label>
				<input
					bind:value={alias}
					type="text"
					id="alias"
					name="alias"
					class="input input-bordered"
					placeholder="New alias..."
				/>
			</div>

			<div class="px-2 py-3">@</div>

			<div>
				<label for="domain" class="sr-only">Domain</label>
				<select
					id="domain"
					name="domain"
					class="select select-bordered"
					bind:value={domain}
				>
					<option value="">Select domain...</option>
					{#each domainOptions as option}
						<option value={option}>{option}</option>
					{/each}
				</select>
			</div>
		</div>

		<div class="add-row">
			<p class="text-md text-primary">Redirects to</p>
		</div>

		<div class="add-row flex items-center">
			<div>
				<label for="email" class="sr-only">Email</label>
				<select
					id="email"
					name="email"
					class="select select-bordered"
					bind:value={email}
				>
					<option value="">Select email...</option>
					{#each emailSelectOptions as option}
						<option value={option}>{option}</option>
					{/each}
				</select>
			</div>

			<div>
				<label class="pl-2 label cursor-pointer">
					<input
						id="includeExistingAliases"
						name="includeExistingAliases"
						type="checkbox"
						bind:checked={includeExistingAliases}
						class="checkbox"
					/>
					<span class="pl-2 label-text">Include aliases</span>
				</label>
			</div>
		</div>

		{#if isLoading}
			<Spinner />
		{:else}
			<div class="add-row">
				<button
					type="submit"
					class="btn btn-primary"
					disabled={validAlias !== true}
				>
					Add
				</button>
			</div>
		{/if}
	</form>
</div>

<style>
	@reference "../app.css";
	.add-row {
		@apply flex justify-center mb-4;
	}
</style>
