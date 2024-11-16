<script lang='ts'>
  import DestinationCard from '$lib/destinationCard.svelte';

  let tmpImage = 'rocks.jpg';
</script>

<DestinationCard img={tmpImage} endDate={new Date(Date.now())} name="Rocks" />
