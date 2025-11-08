import type { FC, FormEvent } from "react";
import { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom'
import { logoutStaff } from '../lib/auth'
import { apiService, type EntityType, type GenericItem } from '../lib/api'

const emptyTemplate = (tipo: EntityType): GenericItem => {
	const base: GenericItem = { id: "", nombre: "" };
	switch (tipo) {
		case "software":
			return { ...base, tipo: "", empresa: "", image_url: "", orden_grafica: 0, orden_procesador: 0, orden_ram: 0, carrera: "" };
		case "gabinete":
			return { ...base, fabricante: "", form_factor: "", image_url: "", max_largo_gpu_float: 0, msrp: "", socket: "" };
		case "procesador":
			return { ...base, consumo: 0, fabricante: "", generacion: "", image_url: "", msrp: "", nucleos: 0 };
		case "placaGrafica":
			return { ...base, consumo: 0, fabricante: "", image_url: "", largo: 0, modelo: "", vram: 0 };
		case "almacenamiento":
			return { ...base, capacidad: 0, fabricante: "", image_url: "", msrp: "", tipo_almacenamiento: "" };
		case "placaMadre":
			return { ...base, fabricante: "", form_factor: "", image_url: "", modelo: "", msrp: "", socket_procesador: "", socket_ram: "" };
		case "fuente":
			return { ...base, capacidad: 0, fabricante: "", image_url: "" };
		case "memoriaRam":
			return { ...base, capacidad: 0, fabricante: "", generacion: "", image_url: "" };
		case "carrera":
			return { ...base };
	}
};

const PanelDeAdmin: FC = () => {
	const [items, setItems] = useState<GenericItem[]>([]);
	const [selectedType, setSelectedType] = useState<EntityType>('software');
	const [form, setForm] = useState<GenericItem>(() => emptyTemplate('software'));
	const [editingId, setEditingId] = useState<string | null>(null);
	const [viewMode, setViewMode] = useState<'choose' | 'form' | 'list'>('choose');
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [statsData, setStatsData] = useState<Record<EntityType, number>>({} as Record<EntityType, number>);
	const [statsLoading, setStatsLoading] = useState(true);
	const [searchTerm, setSearchTerm] = useState<string>('');
	const [filteredSearchTerm, setFilteredSearchTerm] = useState<string>('');
	const [carreras, setCarreras] = useState<GenericItem[]>([]);
	const [carrerasLoading, setCarrerasLoading] = useState(true);
	const [selectedImage, setSelectedImage] = useState<File | null>(null);
	const [imagePreview, setImagePreview] = useState<string | null>(null);
	const [dragActive, setDragActive] = useState(false);
	const navigate = useNavigate();

	useEffect(() => {
		if (viewMode === 'list') {
			loadItems();
		}
		// Limpiar búsqueda e imagen al cambiar de vista
		setSearchTerm('');
		setFilteredSearchTerm('');
		setSelectedImage(null);
		setImagePreview(null);
	}, [selectedType, viewMode]);

	useEffect(() => {
		loadStatsData();
		loadCarreras();
	}, []);

	async function loadCarreras() {
		setCarrerasLoading(true);
		try {
			const data = await apiService.list('carrera');
			setCarreras(data);
		} catch (err) {
			console.error('Error loading carreras:', err);
		} finally {
			setCarrerasLoading(false);
		}
	}

	function performSearch() {
		setFilteredSearchTerm(searchTerm.trim());
	}

	// Funciones para manejo de imágenes
	function validateImageFile(file: File): boolean {
		// Validar tamaño (10MB máximo)
		const maxSize = 10 * 1024 * 1024; // 10MB en bytes
		if (file.size > maxSize) {
			setError('La imagen no puede ser mayor a 10MB');
			return false;
		}

		// Validar tipo de archivo (solo JPG)
		const validTypes = ['image/jpeg', 'image/jpg'];
		if (!validTypes.includes(file.type)) {
			setError('Solo se permiten archivos JPG');
			return false;
		}

		return true;
	}

	function handleImageSelect(file: File) {
		if (validateImageFile(file)) {
			setSelectedImage(file);
			
			// Crear preview
			const reader = new FileReader();
			reader.onload = (e) => {
				setImagePreview(e.target?.result as string);
			};
			reader.readAsDataURL(file);
			setError(null);
		}
	}

	function handleImageDrop(e: React.DragEvent) {
		e.preventDefault();
		setDragActive(false);
		
		const files = e.dataTransfer.files;
		if (files.length > 0) {
			handleImageSelect(files[0]);
		}
	}

	function handleImageDragOver(e: React.DragEvent) {
		e.preventDefault();
		setDragActive(true);
	}

	function handleImageDragLeave(e: React.DragEvent) {
		e.preventDefault();
		setDragActive(false);
	}

	function removeImage() {
		setSelectedImage(null);
		setImagePreview(null);
	}

	async function loadStatsData() {
		setStatsLoading(true);
		const stats: Record<EntityType, number> = {} as Record<EntityType, number>;
		
		try {
			const entityTypes: EntityType[] = [
				'software', 'carrera', 'gabinete', 'procesador', 
				'placaGrafica', 'almacenamiento', 'placaMadre', 'fuente', 'memoriaRam'
			];
			
			const promises = entityTypes.map(async (type) => {
				try {
					const data = await apiService.list(type);
					stats[type] = data.length;
				} catch (err) {
					console.error(`Error loading ${type}:`, err);
					stats[type] = 0;
				}
			});
			
			await Promise.all(promises);
			setStatsData(stats);
		} catch (err) {
			console.error('Error loading stats:', err);
		} finally {
			setStatsLoading(false);
		}
	}

	async function loadItems() {
		setLoading(true);
		setError(null);
		try {
			const data = await apiService.list(selectedType);
			setItems(data);
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Error al cargar items');
		} finally {
			setLoading(false);
		}
	}

	function handleChange(key: string, value: any) {
		setForm((f) => ({ ...f, [key]: value }));
	}

	async function onSubmit(e: FormEvent) {
		e.preventDefault();
		if (!form.nombre?.trim()) return;

		setLoading(true);
		setError(null);
		
		try {
			if (editingId) {
				const updatedItem = await apiService.update(selectedType, editingId, form, selectedImage || undefined);
				setItems((prev) => prev.map((it) => (it.id === editingId ? updatedItem : it)));
				setEditingId(null);
			} else {
				const newItem = await apiService.create(selectedType, form, selectedImage || undefined);
				setItems((prev) => [newItem, ...prev]);
				// Actualizar estadísticas después de crear un nuevo item
				setStatsData(prev => ({
					...prev,
					[selectedType]: (prev[selectedType] || 0) + 1
				}));
			}
			setForm(emptyTemplate(selectedType));
			setSelectedImage(null);
			setImagePreview(null);
			setViewMode('list');
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Error al guardar item');
		} finally {
			setLoading(false);
		}
	}

	function onEdit(id: string) {
		const it = items.find((x) => x.id === id);
		if (!it) return;
		setForm(it);
		setEditingId(id);
		setViewMode('form');
	}

	async function onDelete(id: string) {
		if (!confirm("¿Eliminar este item?")) return;
		
		setLoading(true);
		setError(null);
		
		try {
			await apiService.delete(selectedType, id);
			setItems((prev) => prev.filter((x) => x.id !== id));
			// Actualizar estadísticas después de eliminar un item
			setStatsData(prev => ({
				...prev,
				[selectedType]: Math.max(0, (prev[selectedType] || 0) - 1)
			}));
			if (editingId === id) {
				setEditingId(null);
				setForm(emptyTemplate(selectedType));
			}
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Error al eliminar item');
		} finally {
			setLoading(false);
		}
	}

	// Filtrar items según el término de búsqueda (optimizado por tipo de componente)
	const filteredItems = items.filter(item => {
		if (!filteredSearchTerm) return true;
		
		const searchLower = filteredSearchTerm.toLowerCase();
		const matchesName = item.nombre?.toLowerCase().includes(searchLower);
		
		// Búsqueda específica por tipo de componente
		switch (selectedType) {
			case 'software':
				const matchesTipo = item.tipo?.toLowerCase().includes(searchLower);
				const matchesEmpresa = item.empresa?.toLowerCase().includes(searchLower);
				const matchesCarrera = item.carrera?.toLowerCase().includes(searchLower);
				return matchesName || matchesTipo || matchesEmpresa || matchesCarrera;
				
			case 'procesador':
				const matchesFabricanteCPU = item.fabricante?.toLowerCase().includes(searchLower);
				const matchesGeneracion = item.generacion?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricanteCPU || matchesGeneracion;
				
			case 'placaGrafica':
				const matchesFabricanteGPU = item.fabricante?.toLowerCase().includes(searchLower);
				const matchesModeloGPU = item.modelo?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricanteGPU || matchesModeloGPU;
				
			case 'placaMadre':
				const matchesModeloMB = item.modelo?.toLowerCase().includes(searchLower);
				const matchesFabricanteMB = item.fabricante?.toLowerCase().includes(searchLower);
				return matchesName || matchesModeloMB || matchesFabricanteMB;
				
			case 'almacenamiento':
				const matchesFabricanteHDD = item.fabricante?.toLowerCase().includes(searchLower);
				const matchesTipoAlmacenamiento = item.tipo_almacenamiento?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricanteHDD || matchesTipoAlmacenamiento;
				
			case 'memoriaRam':
				const matchesFabricanteRAM = item.fabricante?.toLowerCase().includes(searchLower);
				const matchesGeneracionRAM = item.generacion?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricanteRAM || matchesGeneracionRAM;
				
			case 'gabinete':
				const matchesFabricanteCase = item.fabricante?.toLowerCase().includes(searchLower);
				const matchesFormFactor = item.form_factor?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricanteCase || matchesFormFactor;
				
			case 'fuente':
				const matchesFabricantePSU = item.fabricante?.toLowerCase().includes(searchLower);
				return matchesName || matchesFabricantePSU;
				
			case 'carrera':
				// Solo buscar por nombre para carreras
				return matchesName;
				
			default:
				return matchesName;
		}
	});

	// Obtener placeholder dinámico según el tipo de componente
	const getSearchPlaceholder = (type: EntityType): string => {
		switch (type) {
			case 'software':
				return 'Buscar software por nombre, tipo, empresa o carrera...';
			case 'procesador':
				return 'Buscar procesador por nombre, fabricante o generación...';
			case 'placaGrafica':
				return 'Buscar placa gráfica por nombre, fabricante o modelo...';
			case 'placaMadre':
				return 'Buscar placa madre por nombre, modelo o fabricante...';
			case 'almacenamiento':
				return 'Buscar almacenamiento por nombre, fabricante o tipo...';
			case 'memoriaRam':
				return 'Buscar memoria RAM por nombre, fabricante o generación...';
			case 'gabinete':
				return 'Buscar gabinete por nombre, fabricante o form factor...';
			case 'fuente':
				return 'Buscar fuente por nombre o fabricante...';
			case 'carrera':
				return 'Buscar carrera por nombre...';
			default:
				return `Buscar ${type} por nombre...`;
		}
	};

	const entities = [
		{ key: 'software', label: 'Software', icon: 'apps', description: 'Aplicaciones y juegos' },
		{ key: 'carrera', label: 'Carreras', icon: 'school', description: 'Carreras universitarias' },
		{ key: 'gabinete', label: 'Gabinetes', icon: 'inventory_2', description: 'Cases y chasis' },
		{ key: 'procesador', label: 'Procesadores', icon: 'memory', description: 'CPU y chips' },
		{ key: 'placaGrafica', label: 'Placas Gráficas', icon: 'developer_board', description: 'GPU y tarjetas gráficas' },
		{ key: 'almacenamiento', label: 'Almacenamiento', icon: 'save', description: 'SSD, HDD y discos' },
		{ key: 'placaMadre', label: 'Placas Madre', icon: 'account_tree', description: 'Motherboards' },
		{ key: 'fuente', label: 'Fuentes', icon: 'power', description: 'PSU y alimentación' },
		{ key: 'memoriaRam', label: 'Memoria RAM', icon: 'storage', description: 'Memoria del sistema' },
	] as { key: EntityType; label: string; icon: string; description: string }[];

	const renderFormFields = () => {
		const inputClass = "w-full px-4 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent";
		const labelClass = "block text-sm font-medium text-gray-600 dark:text-gray-400 mb-2";
		
		switch (selectedType) {
			case 'software':
				return (
					<div className="space-y-4">
						<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div>
								<label className={labelClass}>Tipo</label>
								<input value={form.tipo || ''} onChange={(e) => handleChange('tipo', e.target.value)} className={inputClass} />
							</div>
							<div>
								<label className={labelClass}>Empresa</label>
								<input value={form.empresa || ''} onChange={(e) => handleChange('empresa', e.target.value)} className={inputClass} />
							</div>
						</div>
						<div>
							<label className={labelClass}>Image URL (opcional)</label>
							<input value={form.image_url || ''} onChange={(e) => handleChange('image_url', e.target.value)} className={inputClass} placeholder="https://ejemplo.com/imagen.jpg" />
						</div>
						
						{/* Componente de carga de imagen */}
						<div>
							<label className={labelClass}>Subir imagen (opcional)</label>
							<div className="space-y-3">
								<div
									className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
										dragActive 
											? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' 
											: 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
									}`}
									onDrop={handleImageDrop}
									onDragOver={handleImageDragOver}
									onDragLeave={handleImageDragLeave}
								>
									{imagePreview ? (
										<div className="space-y-3">
											<img 
												src={imagePreview} 
												alt="Preview" 
												className="mx-auto max-h-32 rounded-lg object-contain"
											/>
											<div className="flex justify-center gap-2">
												<label className="cursor-pointer px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors">
													Cambiar imagen
													<input
														type="file"
														accept="image/jpeg,image/jpg"
														onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0])}
														className="hidden"
													/>
												</label>
												<button
													type="button"
													onClick={removeImage}
													className="px-3 py-1 bg-red-600 hover:bg-red-700 text-white rounded text-sm transition-colors"
												>
													Eliminar
												</button>
											</div>
										</div>
									) : (
										<div className="space-y-3">
											<div className="mx-auto w-12 h-12 text-gray-400">
												<span className="material-symbols-outlined text-4xl">cloud_upload</span>
											</div>
											<div>
												<p className="text-gray-600 dark:text-gray-400 mb-2">
													Arrastra una imagen aquí o
												</p>
												<label className="cursor-pointer px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors inline-block">
													Seleccionar archivo
													<input
														type="file"
														accept="image/jpeg,image/jpg"
														onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0])}
														className="hidden"
													/>
												</label>
											</div>
											<p className="text-xs text-gray-500 dark:text-gray-400">
												JPG, PNG, GIF, WebP (máx. 10MB)
											</p>
										</div>
									)}
								</div>
								{selectedImage && (
									<div className="text-sm text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 rounded p-2">
										<strong>Archivo seleccionado:</strong> {selectedImage.name} ({(selectedImage.size / 1024 / 1024).toFixed(2)} MB)
									</div>
								)}
							</div>
						</div>
						
						<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
							<div>
								<label className={labelClass}>Orden gráfica</label>
								<input type="number" value={form.orden_grafica || 0} onChange={(e) => handleChange('orden_grafica', Number(e.target.value))} className={inputClass} />
							</div>
							<div>
								<label className={labelClass}>Orden procesador</label>
								<input type="number" value={form.orden_procesador || 0} onChange={(e) => handleChange('orden_procesador', Number(e.target.value))} className={inputClass} />
							</div>
							<div>
								<label className={labelClass}>Orden RAM</label>
								<input type="number" value={form.orden_ram || 0} onChange={(e) => handleChange('orden_ram', Number(e.target.value))} className={inputClass} />
							</div>
						</div>
						<div>
							<label className={labelClass}>Carrera *</label>
							<select 
								value={form.carrera || ''} 
								onChange={(e) => handleChange('carrera', e.target.value)} 
								className={inputClass}
								required
								disabled={carrerasLoading}
							>
								<option value="">
									{carrerasLoading ? 'Cargando carreras...' : 'Seleccionar carrera...'}
								</option>
								<option value="videojuego">Video juego</option>
								{carreras.length > 0 ? (
									carreras.map((carrera) => (
										<option key={carrera.id} value={carrera.nombre}>
											{carrera.nombre}
										</option>
									))
								) : (
									!carrerasLoading && (
										<option value="" disabled>
											No hay carreras disponibles
										</option>
									)
								)}
							</select>
							{carrerasLoading && (
								<p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
									Cargando carreras desde la base de datos...
								</p>
							)}
							{!carrerasLoading && carreras.length === 0 && (
								<p className="text-xs text-yellow-600 dark:text-yellow-400 mt-1">
									⚠️ No se encontraron carreras en la base de datos
								</p>
							)}
							{!carrerasLoading && carreras.length > 0 && (
								<p className="text-xs text-green-600 dark:text-green-400 mt-1">
									✅ {carreras.length} carrera(s) cargada(s) desde la base de datos
								</p>
							)}
						</div>
					</div>
				);
			case 'carrera':
				return <p className="text-gray-600 dark:text-gray-400">Solo es necesario completar el nombre para las carreras.</p>;
			default:
				return (
					<div className="space-y-4">
						<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div>
								<label className={labelClass}>Fabricante</label>
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} className={inputClass} />
							</div>
							{form.modelo !== undefined && (
								<div>
									<label className={labelClass}>Modelo</label>
									<input value={form.modelo || ''} onChange={(e) => handleChange('modelo', e.target.value)} className={inputClass} />
								</div>
							)}
						</div>
						<div>
							<label className={labelClass}>Image URL (opcional)</label>
							<input value={form.image_url || ''} onChange={(e) => handleChange('image_url', e.target.value)} className={inputClass} placeholder="https://ejemplo.com/imagen.jpg" />
						</div>
						
						{/* Componente de carga de imagen */}
						<div>
							<label className={labelClass}>Subir imagen (opcional)</label>
							<div className="space-y-3">
								<div
									className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
										dragActive 
											? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' 
											: 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
									}`}
									onDrop={handleImageDrop}
									onDragOver={handleImageDragOver}
									onDragLeave={handleImageDragLeave}
								>
									{imagePreview ? (
										<div className="space-y-3">
											<img 
												src={imagePreview} 
												alt="Preview" 
												className="mx-auto max-h-32 rounded-lg object-contain"
											/>
											<div className="flex justify-center gap-2">
												<label className="cursor-pointer px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors">
													Cambiar imagen
													<input
														type="file"
														accept="image/jpeg,image/jpg"
														onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0])}
														className="hidden"
													/>
												</label>
												<button
													type="button"
													onClick={removeImage}
													className="px-3 py-1 bg-red-600 hover:bg-red-700 text-white rounded text-sm transition-colors"
												>
													Eliminar
												</button>
											</div>
										</div>
									) : (
										<div className="space-y-3">
											<div className="mx-auto w-12 h-12 text-gray-400">
												<span className="material-symbols-outlined text-4xl">cloud_upload</span>
											</div>
											<div>
												<p className="text-gray-600 dark:text-gray-400 mb-2">
													Arrastra una imagen aquí o
												</p>
												<label className="cursor-pointer px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors inline-block">
													Seleccionar archivo
													<input
														type="file"
														accept="image/jpeg,image/jpg"
														onChange={(e) => e.target.files?.[0] && handleImageSelect(e.target.files[0])}
														className="hidden"
													/>
												</label>
											</div>
											<p className="text-xs text-gray-500 dark:text-gray-400">
												JPG, PNG, GIF, WebP (máx. 10MB)
											</p>
										</div>
									)}
								</div>
								{selectedImage && (
									<div className="text-sm text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 rounded p-2">
										<strong>Archivo seleccionado:</strong> {selectedImage.name} ({(selectedImage.size / 1024 / 1024).toFixed(2)} MB)
									</div>
								)}
							</div>
						</div>
						
						{form.msrp !== undefined && (
							<div>
								<label className={labelClass}>MSRP ($)</label>
								<input type="text" value={form.msrp || ''} onChange={(e) => handleChange('msrp', e.target.value)} className={inputClass} placeholder="Ej: $299.99" />
							</div>
						)}
					</div>
				);
		}
	};

	return (
		<div className="relative flex min-h-screen bg-gray-50 dark:bg-gray-900 font-sans">
			{/* Sidebar Navigation */}
			<aside className="w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col">
				<div className="p-4 border-b border-gray-200 dark:border-gray-700">
					<h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">PC Recs Admin</h1>
				</div>
				<nav className="flex-1 p-4 space-y-2">
					<a 
						href="#" 
						className={`flex items-center gap-3 px-3 py-2 rounded-lg font-medium ${
							viewMode === 'choose' 
								? 'bg-blue-500/20 text-blue-600 dark:text-blue-400' 
								: 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
						}`}
						onClick={(e) => { e.preventDefault(); setViewMode('choose'); }}
					>
						<span className="material-symbols-outlined text-xl">dashboard</span>
						<span>Dashboard</span>
					</a>
					{entities.map((entity) => (
						<a
							key={entity.key}
							href="#"
							className={`flex items-center gap-3 px-3 py-2 rounded-lg font-medium ${
								selectedType === entity.key && viewMode !== 'choose'
									? 'bg-blue-500/20 text-blue-600 dark:text-blue-400'
									: 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
							}`}
							onClick={(e) => {
								e.preventDefault();
								setSelectedType(entity.key);
								setForm(emptyTemplate(entity.key));
								setEditingId(null);
								setViewMode('list');
							}}
						>
							<span className="material-symbols-outlined text-xl">{entity.icon}</span>
							<span>{entity.label}</span>
						</a>
					))}
				</nav>
				<div className="p-4 border-t border-gray-200 dark:border-gray-700">
					<button
						onClick={() => { logoutStaff(); navigate('/login'); }}
						className="flex items-center gap-3 px-3 py-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 w-full font-medium"
					>
						<span className="material-symbols-outlined text-xl">logout</span>
						<span>Cerrar sesión</span>
					</button>
				</div>
			</aside>

			{/* Main Content */}
			<div className="flex-1 flex flex-col">
				{/* Top Header Bar */}
				<header className="flex items-center justify-between p-4 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
					<div className="flex items-center gap-4 flex-1">
						{viewMode === 'list' && (
							<div className="relative w-full max-w-md">
								<span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">search</span>
								<input 
									className="w-full pl-10 pr-16 py-2 rounded-lg bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 dark:text-gray-100" 
									placeholder={getSearchPlaceholder(selectedType)}
									type="text"
									value={searchTerm}
									onChange={(e) => setSearchTerm(e.target.value)}
									onKeyDown={(e) => e.key === 'Enter' && performSearch()}
								/>
								<button
									className="absolute right-3 top-1/2 -translate-y-1/2 text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
									onClick={performSearch}
								>
									<span className="material-symbols-outlined text-sm">search</span>
								</button>
							</div>
						)}
						{viewMode === 'choose' && (
							<button className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium transition-colors whitespace-nowrap">
								<span className="material-symbols-outlined text-sm">sync</span>
								Scrapear datos
							</button>
						)}
					</div>
					<div className="flex items-center gap-4">
						<button className="p-2 rounded-full text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700">
							<span className="material-symbols-outlined">light_mode</span>
						</button>
						<div className="relative">
							<button className="flex items-center gap-2">
								<div className="bg-center bg-no-repeat aspect-square bg-cover rounded-full size-10 bg-gray-300 dark:bg-gray-600"></div>
							</button>
						</div>
					</div>
				</header>

				{/* Main Content Area */}
				<main className="flex-1 p-6 space-y-6 bg-gray-50 dark:bg-gray-900">
					{viewMode === 'choose' && (
						<div>
							{/* Page Heading */}
							<div className="flex flex-wrap justify-between gap-3 mb-6">
								<p className="text-gray-900 dark:text-gray-100 text-4xl font-black leading-tight tracking-[-0.033em] min-w-72">¡Bienvenido, Administrador!</p>
							</div>

							{/* Stats */}
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
									<p className="text-gray-600 dark:text-gray-400 text-base font-medium leading-normal">Total Software</p>
									<p className="text-gray-900 dark:text-gray-100 tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.software || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
									<p className="text-gray-600 dark:text-gray-400 text-base font-medium leading-normal">Total Procesadores</p>
									<p className="text-gray-900 dark:text-gray-100 tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.procesador || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
									<p className="text-gray-600 dark:text-gray-400 text-base font-medium leading-normal">Total Placas Gráficas</p>
									<p className="text-gray-900 dark:text-gray-100 tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.placaGrafica || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
									<p className="text-gray-600 dark:text-gray-400 text-base font-medium leading-normal">Total Almacenamiento</p>
									<p className="text-gray-900 dark:text-gray-100 tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.almacenamiento || 0)}
									</p>
								</div>
							</div>

							{/* Complete Stats Table */}
							<div className="mb-8">
								<div className="flex justify-between items-center pb-3">
									<h2 className="text-gray-900 dark:text-gray-100 text-[22px] font-bold leading-tight tracking-[-0.015em]">Resumen Completo</h2>
								</div>
								<div className="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
									<div className="overflow-x-auto">
										<table className="w-full">
											<thead>
												<tr className="bg-gray-50 dark:bg-gray-700">
													<th className="px-6 py-3 text-left text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Categoría</th>
													<th className="px-6 py-3 text-left text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Descripción</th>
													<th className="px-6 py-3 text-center text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Cantidad</th>
													<th className="px-6 py-3 text-center text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Estado</th>
												</tr>
											</thead>
											<tbody className="divide-y divide-gray-200 dark:divide-gray-700">
												{entities.map((entity) => {
													const count = statsLoading ? '...' : (statsData[entity.key] || 0);
													const hasItems = !statsLoading && (statsData[entity.key] || 0) > 0;
													
													return (
														<tr key={entity.key} className="hover:bg-gray-50 dark:hover:bg-gray-700">
															<td className="px-6 py-4 whitespace-nowrap">
																<div className="flex items-center gap-3">
																	<span className="material-symbols-outlined text-2xl text-blue-600 dark:text-blue-400">{entity.icon}</span>
																	<div className="text-sm font-medium text-gray-900 dark:text-gray-100">{entity.label}</div>
																</div>
															</td>
															<td className="px-6 py-4 whitespace-nowrap">
																<div className="text-sm text-gray-600 dark:text-gray-400">{entity.description}</div>
															</td>
															<td className="px-6 py-4 whitespace-nowrap text-center">
																<div className="text-lg font-bold text-gray-900 dark:text-gray-100">{count}</div>
															</td>
															<td className="px-6 py-4 whitespace-nowrap text-center">
																{statsLoading ? (
																	<span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300">
																		Cargando...
																	</span>
																) : hasItems ? (
																	<span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300">
																		✅ Con datos
																	</span>
																) : (
																	<span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300">
																		⚠️ Sin datos
																	</span>
																)}
															</td>
														</tr>
													);
												})}
											</tbody>
										</table>
									</div>
								</div>
							</div>
						</div>
					)}

					{viewMode === 'form' && (
						<div className="max-w-2xl">
							<div className="mb-6">
								<button 
									className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-medium transition-colors"
									onClick={() => setViewMode('list')}
								>
									<span className="material-symbols-outlined">arrow_back</span>
									Volver a la lista
								</button>
							</div>

							<div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
								<div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
									<h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
										{editingId ? `Editar ${selectedType}` : `Nuevo ${selectedType}`}
									</h2>
								</div>
								
								<form className="p-6 space-y-6" onSubmit={onSubmit}>
									<div>
										<label className="block text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
											Nombre *
										</label>
										<input 
											value={form.nombre} 
											onChange={(e) => handleChange('nombre', e.target.value)} 
											required 
											className="w-full px-4 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
											placeholder="Ingrese el nombre"
										/>
									</div>

									{renderFormFields()}

									{error && (
										<div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg">
											<strong>Error:</strong> {error}
										</div>
									)}

									<div>
										<button 
											type="submit" 
											disabled={loading}
											className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white rounded-lg font-medium transition-colors"
										>
											{loading ? 'Guardando...' : (editingId ? "Guardar" : "Crear")}
										</button>
									</div>
								</form>
							</div>
						</div>
					)}

					{viewMode === 'list' && (
						<div>
							<div className="mb-6">
								<button 
									className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-medium transition-colors"
									onClick={() => setViewMode('choose')}
								>
									<span className="material-symbols-outlined">arrow_back</span>
									Volver
								</button>
							</div>

							<div className="flex justify-between items-center mb-6">
								<div>
									<h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 capitalize">{selectedType}</h2>
									<p className="text-gray-600 dark:text-gray-400">
										{filteredSearchTerm ? 
											`${filteredItems.length} de ${items.length} items (buscando: "${filteredSearchTerm}")` : 
											`Total de items: ${items.length}`
										}
									</p>
								</div>
								<button 
									className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
									onClick={() => setViewMode('form')}
								>
									<span className="material-symbols-outlined">add</span>
									Agregar nuevo
								</button>
							</div>

							{error && (
								<div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg mb-4">
									<strong>Error:</strong> {error}
								</div>
							)}

							{loading ? (
								<div className="text-center py-8">
									<div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
									<p className="mt-2 text-gray-600 dark:text-gray-400">Cargando...</p>
								</div>
							) : items.length === 0 ? (
								<div className="text-center py-12">
									<p className="text-gray-500 dark:text-gray-400 text-lg">No hay items cargados aún</p>
								</div>
							) : filteredItems.length === 0 ? (
								<div className="text-center py-12">
									<p className="text-gray-500 dark:text-gray-400 text-lg">No se encontraron resultados para "{filteredSearchTerm}"</p>
									<button
										className="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
										onClick={() => {
											setSearchTerm('');
											setFilteredSearchTerm('');
										}}
									>
										Limpiar búsqueda
									</button>
								</div>
							) : (
								<div className="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
									<div className="overflow-x-auto">
										<table className="w-full">
											<thead>
												<tr className="bg-gray-50 dark:bg-gray-700">
													<th className="px-6 py-3 text-left text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Nombre</th>
													<th className="px-6 py-3 text-left text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Detalles</th>
													<th className="px-6 py-3 text-left text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Specs</th>
													<th className="px-6 py-3 text-right text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider">Acciones</th>
												</tr>
											</thead>
											<tbody className="divide-y divide-gray-200 dark:divide-gray-700">
												{filteredItems.map((it) => (
													<tr key={it.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
														<td className="px-6 py-4 whitespace-nowrap">
															<div className="text-sm font-medium text-gray-900 dark:text-gray-100">{it.nombre}</div>
														</td>
														<td className="px-6 py-4 whitespace-nowrap">
															<div className="text-sm text-gray-600 dark:text-gray-400">
																{it.fabricante || it.tipo || ''} {it.modelo || it.empresa ? `— ${it.modelo || it.empresa}` : ''}
															</div>
														</td>
														<td className="px-6 py-4 whitespace-nowrap">
															<div className="text-sm text-gray-600 dark:text-gray-400">
																{it.vram ? `VRAM:${it.vram}GB ` : ''}
																{it.nucleos ? `Núcleos:${it.nucleos} ` : ''}
																{it.capacidad ? `Cap:${it.capacidad} ` : ''}
																{it.msrp ? `$${it.msrp}` : ''}
															</div>
														</td>
														<td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
															<button 
																className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 mr-4 disabled:opacity-50"
																onClick={() => onEdit(it.id)} 
																disabled={loading}
															>
																<span className="material-symbols-outlined text-base">edit</span>
															</button>
															<button 
																className="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 disabled:opacity-50"
																onClick={() => onDelete(it.id)} 
																disabled={loading}
															>
																<span className="material-symbols-outlined text-base">delete</span>
															</button>
														</td>
													</tr>
												))}
											</tbody>
										</table>
									</div>
								</div>
							)}
						</div>
					)}
				</main>
			</div>
		</div>
	);
};

export default PanelDeAdmin;
