// Servicio para interactuar con la API del backend
// Lee la URL base desde variables de entorno de Vite (prefijo VITE_),
// con fallback a misma origin + "/v1". En desarrollo configurar en .env.development.
import { getApiBase } from './env';
import { getToken } from './auth';
const API_BASE_URL = getApiBase();

export type EntityType =
  | "software"
  | "gabinete" 
  | "procesador"
  | "placaGrafica"
  | "almacenamiento"
  | "placaMadre"
  | "fuente"
  | "memoriaRam"
  | "carrera";

// Mapeo de tipos del frontend a rutas de la API
const ENTITY_ENDPOINTS: Record<EntityType, string> = {
  software: 'softwares',
  gabinete: 'componentes/gabinetes', 
  procesador: 'componentes/procesadores',
  placaGrafica: 'componentes/gpus',
  almacenamiento: 'componentes/discos',
  placaMadre: 'componentes/motherboards',
  fuente: 'componentes/psus',
  memoriaRam: 'componentes/rams',
  carrera: 'carreras'
};

export interface GenericItem {
  id: string;
  nombre: string;
  
  // Software fields
  tipo?: string;
  empresa?: string;
  image_url?: string;
  orden_grafica?: number;
  orden_procesador?: number;
  orden_ram?: number;
  carrera?: string;
  
  // Component fields
  fabricante?: string;
  msrp?: string;
  
  // Gabinete
  form_factor?: string;
  max_largo_gpu_float?: number;
  socket?: string;
  
  // Procesador
  consumo?: number;
  generacion?: string;
  nucleos?: number;
  
  // GPU
  largo?: number;
  modelo?: string;
  vram?: number;
  
  // Disco
  capacidad?: number;
  tipo_almacenamiento?: string;
  
  // Motherboard  
  socket_procesador?: string;
  socket_ram?: string;
  
  // RAM
  // capacidad ya definida arriba
  // generacion ya definida arriba
}

class ApiService {
  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}/${endpoint}`;
    
    
    console.log(options.body);

  // Crear headers base
  const headers = new Headers();
  
  // Agregar token si existe
  const token = getToken();
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  
  // Agregar Content-Type si no es FormData
  if (!(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  // Combinar con headers personalizados
  if (options.headers) {
    Object.entries(options.headers).forEach(([key, value]) => {
      headers.set(key, value);
    });
  }

  
  const config: RequestInit = {
    ...options,
    headers,
    credentials: options.credentials,
  };

    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`API Error ${response.status}: ${errorText}`);
    }
    
    // Si no hay contenido devolver undefined para permitir métodos que esperan void
    if (response.status === 204) {
      return undefined as unknown as T;
    }

    // Determinar tipo de contenido
    const contentType = response.headers.get('content-type') || '';

    // Si es JSON, parsear normalmente
    if (contentType.includes('application/json')) {
      return response.json();
    }

    // Si no es JSON, intentar leer como texto. Si está vacío devolver undefined.
    const text = await response.text();
    if (!text) {
      return undefined as unknown as T;
    }

    // Intentar parsear JSON por si el servidor devolvió JSON sin content-type correcto.
    try {
      return JSON.parse(text) as T;
    } catch {
      // Devolver texto crudo si no es JSON
      return text as unknown as T;
    }
  }
  
  // Listar items de una entidad con soporte de query params
  async list(entityType: EntityType, params?: Record<string, string | number | boolean | undefined>): Promise<GenericItem[]> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    let fullEndpoint = endpoint;
    if (params) {
      const search = new URLSearchParams();
      Object.entries(params).forEach(([k, v]) => {
        if (v === undefined || v === null) return;
        search.append(k, String(v));
      });
      const qs = search.toString();
      if (qs.length) fullEndpoint += `?${qs}`;
    }
    return this.request<GenericItem[]>(fullEndpoint);
  }

  // Obtener un item específico
  async get(entityType: EntityType, id: string): Promise<GenericItem> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    return this.request<GenericItem>(`${endpoint}/${id}`);
  }

  // Crear un nuevo item
  async create(entityType: EntityType, data: Partial<GenericItem>, imageFile?: File): Promise<GenericItem> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    
      // Crear FormData
    const formData = new FormData();
    
    // Remover el ID y la imagen URL para creación
    const { id, image_url, ...createData } = data;
    
    // Agregar los datos como campo JSON
    formData.append('data', JSON.stringify(createData));
    
    // Agregar imagen si existe
    if (imageFile) {
      formData.append('imagen', imageFile);
    }

    return this.request<GenericItem>(endpoint, {
      method: 'POST',
      body: formData
    });
  }

  // Actualizar un item existente
  async update(entityType: EntityType, id: string, data: Partial<GenericItem>, imageFile?: File): Promise<GenericItem> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    
    // Crear FormData
    const formData = new FormData();
    
    // Remover el ID y la imagen URL del body
    const { id: _, image_url, ...updateData } = data;
    
    // Agregar los datos como campo JSON
    formData.append('data', JSON.stringify(updateData));
    
    // Agregar imagen si existe
    if (imageFile) {
      formData.append('imagen', imageFile);
    }

    return this.request<GenericItem>(`${endpoint}/${id}`, {
      method: 'PUT',
      body: formData
    });
  }

  // Eliminar un item
  async delete(entityType: EntityType, id: string): Promise<void> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    
    await this.request<void>(`${endpoint}/${id}`, {
      method: 'DELETE',
    });
  }
}

export const apiService = new ApiService();