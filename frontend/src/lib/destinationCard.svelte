<script lang='ts'>
  import { Card } from 'flowbite-svelte';

  let { img, endDate, name }: {img: any, endDate: Date, name: string} = $props();
</script>

<Card {img} class="max-w-[600px]">
  <h3 class="font-bold text-lg text-black">{name}</h3>
  <p class="text-sm text-gray-600">available until {endDate.toISOString().slice(0,10)}</p>
</Card>
