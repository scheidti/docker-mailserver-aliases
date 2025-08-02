<script lang="ts">
	import { baseUrl } from "../config";
	import { toasts } from "../stores";
	import type { AliasResponse } from "../types";
	import ConfirmModal from "./ConfirmModal.svelte";

	const aliasesUrl = baseUrl + "/v1/aliases";

	interface Props {
		aliases?: AliasResponse[];
		refresh?: () => void;
	}

	let { aliases = [], refresh }: Props = $props();
	let isDeleting = $state(false);
	let showModal = $state(false);
	let aliasToDelete = "";

	function isAliasInList(alias: string) {
		return aliases.some((a) => a.alias === alias);
	}

	function confirmDelete(alias: string) {
		aliasToDelete = alias;
		showModal = true;
	}

	async function removeAlias() {
		if (!aliasToDelete || !isAliasInList(aliasToDelete)) {
			return;
		}

		isDeleting = true;

		try {
			const response = await fetch(
				aliasesUrl + "/" + encodeURIComponent(aliasToDelete),
				{
					method: "DELETE",
				},
			);

			if (response.status === 204) {
				refresh?.();
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

		aliasToDelete = "";
		isDeleting = false;
	}
</script>

<div class="overflow-x-auto mx-auto max-w-(--breakpoint-xl)">
	<table class="table">
		<thead>
			<tr>
				<th scope="col">Alias</th>
				<th scope="col">Email</th>
				<th scope="col">Actions</th>
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
							onclick={() => confirmDelete(alias)}
						>
							Delete
						</button>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
<ConfirmModal
	bind:open={showModal}
	title="Delete Alias"
	description="Are you sure you want to delete this alias?"
	confirm={removeAlias}
/>

<style></style>
