<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient } from "casino-sdk";
    import { WS_URL } from "$lib/config";

    onMount(async () => {
        let start = Date.now();
        const client = new CasinoClient(WS_URL);
        await client.connect();
        /*console.log("connected");
        await client.auth.authenticate(
            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWJqZWN0SUQiOjEsImV4cCI6MTc1MDQyNzUwNn0.sjzePA5ozUvJ_FYhAqd6xbQBANQjDpzDx3rvDkHEK80",
            );*/

        //console.log(await client.auth.login("testUser", "test123"));
        console.log(Date.now() - start);
    });
</script>
