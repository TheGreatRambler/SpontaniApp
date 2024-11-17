import type { PageLoad } from './$types';

import type { Destination } from '$lib/destinationInterface.ts';

export const load: PageLoad = async ({ fetch, params }) => {
  let data = {};
  // data.destination = await (await fetch(`${import.meta.env.VITE_BASE_URL}/get?request_type=get_task&task_id=${params.slug}`)).json();
  data.destination = await fetch(`${import.meta.env.VITE_BASE_URL}/get?request_type=get_task&id=${params.slug}`);
  console.log(data.destination);
  const imgObjects = await (await fetch(`${import.meta.env.VITE_BASE_URL}/get?request_type=get_images&task_id=${params.slug}`)).json();
  data.images = imgObjects.map((x) => {return {src: x.url};});
  return data;
};
