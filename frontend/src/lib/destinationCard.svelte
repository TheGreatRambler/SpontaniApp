<script lang='ts'>
  import { Card } from 'flowbite-svelte';

  import SubmitModal from '$lib/submit.svelte';

  let { img, endDate, name }: {img: any, endDate: Date, name: string} = $props();
  let openModal: boolean = $state(false);

  const tzOffset = new Date().getTimezoneOffset() * 60000;
</script>

<Card onclick={() => openModal = true} {img} class="max-w-[600px] hover:bg-gray-200">
  <h3 class="font-bold text-lg text-black">{name.toLowerCase()}</h3>
  <p class="text-sm text-gray-600">available until {new Date(endDate - tzOffset).toISOString().slice(0,10)}</p>
</Card>
<SubmitModal bind:openModal={openModal} {name} />
