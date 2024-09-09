<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { baseUrl } from "../config";
	import { toasts } from "../stores";
	import type { AliasResponse } from "../types";

	const aliasesUrl = baseUrl + "/v1/aliases";
	const dispatch = createEventDispatcher();

	export let aliases: AliasResponse[] = [];
	let isDeleting = false;

	function isAliasInList(alias: string) {
		return aliases.some((a) => a.alias === alias);
	}

	async function removeAlias(alias: string) {
		if (!isAliasInList(alias)) {
			return;
		}

		// TODO: Confirm dialog
		isDeleting = true;

		try {
			const response = await fetch(
				aliasesUrl + "/" + encodeURIComponent(alias),
				{
					method: "DELETE",
				},
			);

			if (response.status === 204) {
				dispatch("refresh");
				toasts.update((toasts) => [
					...toasts,
					{ type: "success", text: "Alias deleted" },
				]);
			} else {
				toasts.update((toasts) => [
					...toasts,
					{
						type: "error",
						text: `Failed to delete alias: ${response.statusText}`,
					},
				]);
			}
		} catch (error) {
			toasts.update((toasts) => [
				...toasts,
				{ type: "error", text: `Failed to delete alias: ${error}` },
			]);
		}

		isDeleting = false;
	}
</script>

<div class="overflow-x-auto mx-auto max-w-screen-xl">
	<table class="table">
		<thead>
			<tr>
				<th scope="col">Alias</th>
				<th scope="col">Email</th>
				<th scope="col"></th>
			</tr>
		</thead>
		<tbody>
			{#each aliases as { alias, email }}
				<tr class="hover">
					<td>{alias}</td>
					<td>{email}</td>
					<td class="w-28">
						<button
							class="btn btn-sm btn-error"
							disabled={isDeleting}
							on:click={() => removeAlias(alias)}
						>
							Delete
						</button>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<style></style>
