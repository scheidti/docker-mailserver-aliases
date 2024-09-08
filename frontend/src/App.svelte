<script lang="ts">
	import AddAlias from "./lib/AddAlias.svelte";
	import Alert from "./lib/Alert.svelte";
	import AliasList from "./lib/AliasList.svelte";
	import Spinner from "./lib/Spinner.svelte";
	import { baseUrl } from "./config";
	import type {
		AliasListResponse,
		AliasResponse,
		StatusResponse,
	} from "./types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const statusUrl = baseUrl + "/v1/status";
	const running = checkIfMailserverIsRunning();

	let isLoading = false;
	let aliases: AliasResponse[] = [];

	async function checkIfMailserverIsRunning() {
		try {
			const response = await fetch(statusUrl);
			const data: StatusResponse = await response.json();
			if (data.running === true) {
				getAliases();
			}
			return data.running === true;
		} catch {
			return false;
		}
	}

	async function getAliases() {
		isLoading = true;

		try {
			const response = await fetch(aliasesUrl);
			const list: AliasListResponse = await response.json();
			aliases = list.aliases;
		} catch {}

		isLoading = false;
	}
</script>

<header>
	<div class="mx-auto flex justify-center p-8">
		<h1 class="font-sans text-secondary drop-shadow-lg text-2xl font-bold">
			Docker Mailserver Aliases
		</h1>
	</div>
</header>

<main>
	{#await running}
		<Spinner />
	{:then isRunning}
		{#if isRunning}
			<AddAlias on:added={getAliases} {aliases} />
			{#if isLoading}
				<div class="flex justify-center">
					<Spinner />
				</div>
			{:else}
				<AliasList on:refresh={getAliases} {aliases} />
			{/if}
		{:else}
			<Alert message={"Mailserver is not running."} type={"error"} />
		{/if}
	{/await}
</main>

<footer>
	<div class="mx-auto flex justify-center p-8 mt-4">
		<div class="text-sm">
			© 2024
			<a
				href="https://github.com/scheidti/docker-mailserver-aliases"
				target="_blank"
				class="link"
			>
				Christian Scheid
			</a>
			- This project is licensed under the MIT license.
		</div>
	</div>
</footer>

<style>
</style>
