const API_URL = process.env.QUASAR_API_URL || 'http://localhost:8083/api/v1'

export type IResponseItem={
  created_at: string,
  id: number,
  type: string,
  user_id: number,
  video_id: number
}

interface IResponse{
  first:boolean,
  items:Array<IResponseItem>,
  last:boolean,
  max_page:number,
  page:number,
  size:number,
  total:number,
  total_pages:number,
  visible:number
}

export const getStatisticByViewType = async (viewType:string) => {
      const raw = await fetchData();
      return raw ? filterData(raw.items, viewType) : null
}

const fetchData = async () : Promise<IResponse | undefined> => {
  let result;

  await fetch(`${API_URL}/stats`, {
    method: 'GET'
  }).then(async (res) => {
    if(res.ok)
      result = await res.json()
  })
  return result
}

const filterData = (items:Array<IResponseItem>, viewType:string) => {
  return items.filter((item) => item.type === viewType)
}
