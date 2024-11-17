<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();
</script>

<p>{data.name}</p>
<p>replace this with the carousel stuff</p>
