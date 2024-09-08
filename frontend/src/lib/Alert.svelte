<script lang="ts">
	export let message: string = "";
	export let type: "error" | "success" | "info" | "warning" = "info";
</script>

<div
	role="alert"
	class={`alert ${type === "error" ? "alert-error" : ""} ${type === "success" ? "alert-success" : ""} ${type === "info" ? "alert-info" : ""} ${type === "warning" ? "alert-warning" : ""}`}
>
	<span>{message}</span>
</div>

<style></style>
