<script lang="ts">
	import { fade } from "svelte/transition";
	import { toasts } from "../stores";

	$: {
		if ($toasts.length > 0) {
			setTimeout(() => {
				toasts.update((toasts) => toasts.slice(1));
			}, 4000);
		}
	}
</script>

<div class="toast toast-top toast-right">
	{#each $toasts as toast}
		<div
			transition:fade={{ delay: 0, duration: 200 }}
			class={`alert ${toast.type === "error" ? "alert-error" : ""} ${toast.type === "success" ? "alert-success" : ""} ${toast.type === "info" ? "alert-info" : ""} ${toast.type === "warning" ? "alert-warning" : ""}`}
		>
			<span>{toast.text}</span>
			<button
				class="btn btn-sm btn-square btn-outline border-0"
				aria-label="Close"
				on:click={() =>
					toasts.update((toasts) => toasts.filter((toast) => toast !== toast))}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</div>
	{/each}
</div>

<style></style>
