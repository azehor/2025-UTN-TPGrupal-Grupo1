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
    
    const headers = new Headers(options.headers || {})
    headers.set('Content-Type', 'application/json')
    const token = getToken()
    if (token) headers.set('Authorization', `Bearer ${token}`)

    const config: RequestInit = {
      ...options,
      headers,
      // We use Bearer tokens, not cookies; avoid forcing credentials to prevent CORS rejections
      // (server currently has AllowCredentials=false). Let callers opt-in if needed.
      credentials: options.credentials,
    };

    // Solo agregar Content-Type si no es FormData
    if (!(options.body instanceof FormData)) {
      config.headers = {
        'Content-Type': 'application/json',
        ...options.headers,
      };
    }

    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`API Error ${response.status}: ${errorText}`);
    }
    
    return response.json();
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
  async create(entityType: EntityType, data: Partial<GenericItem>, _imageFile?: File): Promise<GenericItem> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    
    // Remover el ID y la imagen URL para creación
    const { id, image_url, ...createData } = data;
    
    // TEMPORALMENTE: Solo enviar JSON, ignorar imágenes hasta implementar backend
    return this.request<GenericItem>(endpoint, {
      method: 'POST',
      body: JSON.stringify(createData),
    });
  }

  // Actualizar un item existente
  async update(entityType: EntityType, id: string, data: Partial<GenericItem>, _imageFile?: File): Promise<GenericItem> {
    const endpoint = ENTITY_ENDPOINTS[entityType];
    
    // Remover el ID y la imagen URL del body
    const { id: _, image_url, ...updateData } = data;
    
    // TEMPORALMENTE: Solo enviar JSON, ignorar imágenes hasta implementar backend
    return this.request<GenericItem>(`${endpoint}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updateData),
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