<script lang="ts">
	interface Props {
		message?: string;
		type?: "error" | "success" | "info" | "warning";
	}

	let { message = "", type = "info" }: Props = $props();
</script>

<div
	role="alert"
	class={`alert ${type === "error" ? "alert-error" : ""} ${type === "success" ? "alert-success" : ""} ${type === "info" ? "alert-info" : ""} ${type === "warning" ? "alert-warning" : ""}`}
>
	<span>{message}</span>
</div>

<style></style>
