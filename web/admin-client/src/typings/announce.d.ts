interface AnnounceType {
  id: number;
  title: string;
  content: string;
  createdAt: string;
  url: string;
}

interface AddAnnounceType {
  title: string;
  content: string;
  url: string;
}

interface AnnounceFormType {
  id?: number;
  title: string;
  content: string;
  url: string;
}