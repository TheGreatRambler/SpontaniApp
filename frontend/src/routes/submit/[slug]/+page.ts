import type { PageLoad } from './$types';

import type { Destination } from '$lib/destinationInterface.ts';

export const load: PageLoad = async ({ fetch, params }) => {
  let data = {};
  data.destination = {
    description: 'foo',
    endDate: 0,
    img: '',
    lat: 0.0,
    lng: 0.0,
    name: params.slug
  };
  const imgObjects = await (await fetch(`${import.meta.env.VITE_BASE_URL}/get?request_type=get_images&task_id=${params.slug}`)).json();
  data.images = imgObjects.map((x) => {return {src: x.url};});
  return data;
};
