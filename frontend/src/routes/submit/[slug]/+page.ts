import type { PageLoad } from './$types';

import type { Destination } from '$lib/destinationInterface.ts';

export const load: PageLoad = ({ params }) => {
  const data: Destination = {
    description: 'foo',
    endDate: 0,
    img: '',
    lat: 0.0,
    lng: 0.0,
    name: params.slug
  }
  return data;
};
