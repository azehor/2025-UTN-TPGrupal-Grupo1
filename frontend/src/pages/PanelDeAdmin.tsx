import type { FC, FormEvent } from "react";
import { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom'
import { logoutStaff } from '../lib/auth'
import { apiService, type EntityType, type GenericItem } from '../lib/api'

const emptyTemplate = (tipo: EntityType): GenericItem => {
	const base: GenericItem = { id: "", nombre: "" };
	switch (tipo) {
		case "software":
			return { ...base, tipo: "", empresa: "", image_url: "", orden_grafica: 0, orden_procesador: 0, orden_ram: 0 };
		case "gabinete":
			return { ...base, fabricante: "", form_factor: "", image_url: "", max_largo_gpu_float: 0, msrp: 0, socket: "" };
		case "procesador":
			return { ...base, consumo: 0, fabricante: "", generacion: "", image_url: "", msrp: 0, nucleos: 0 };
		case "placaGrafica":
			return { ...base, consumo: 0, fabricante: "", image_url: "", largo: 0, modelo: "", vram: 0 };
		case "almacenamiento":
			return { ...base, capacidad: 0, fabricante: "", image_url: "", msrp: 0, tipo_almacenamiento: "" };
		case "placaMadre":
			return { ...base, fabricante: "", form_factor: "", image_url: "", modelo: "", msrp: 0, socket_procesador: "", socket_ram: "" };
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
	const navigate = useNavigate();

	useEffect(() => {
		if (viewMode === 'list') {
			loadItems();
		}
	}, [selectedType, viewMode]);

	useEffect(() => {
		loadStatsData();
	}, []);

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
				const updatedItem = await apiService.update(selectedType, editingId, form);
				setItems((prev) => prev.map((it) => (it.id === editingId ? updatedItem : it)));
				setEditingId(null);
			} else {
				const newItem = await apiService.create(selectedType, form);
				setItems((prev) => [newItem, ...prev]);
				// Actualizar estadísticas después de crear un nuevo item
				setStatsData(prev => ({
					...prev,
					[selectedType]: (prev[selectedType] || 0) + 1
				}));
			}
			setForm(emptyTemplate(selectedType));
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
							<label className={labelClass}>Image URL</label>
							<input value={form.image_url || ''} onChange={(e) => handleChange('image_url', e.target.value)} className={inputClass} />
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
							<label className={labelClass}>Image URL</label>
							<input value={form.image_url || ''} onChange={(e) => handleChange('image_url', e.target.value)} className={inputClass} />
						</div>
						{form.msrp !== undefined && (
							<div>
								<label className={labelClass}>MSRP ($)</label>
								<input type="number" step="0.01" value={form.msrp || 0} onChange={(e) => handleChange('msrp', Number(e.target.value))} className={inputClass} />
							</div>
						)}
					</div>
				);
		}
	};

	return (
		<div className="relative flex min-h-screen bg-[#101c22] font-['Space_Grotesk',sans-serif]">
			{/* Sidebar Navigation */}
			<aside className="w-64 bg-[#1a2831] border-r border-gray-600 flex flex-col">
				<div className="p-4 border-b border-gray-600">
					<h1 className="text-2xl font-bold text-white">PC Recs Admin</h1>
				</div>
				<nav className="flex-1 p-4 space-y-2">
					<a 
						href="#" 
						className={`flex items-center gap-3 px-3 py-2 rounded-lg font-medium ${
							viewMode === 'choose' 
								? 'bg-[#13a4ec]/20 text-[#13a4ec]' 
								: 'text-gray-400 hover:bg-gray-700 hover:text-white'
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
									? 'bg-[#13a4ec]/20 text-[#13a4ec]'
									: 'text-gray-400 hover:bg-gray-700 hover:text-white'
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
				<div className="p-4 border-t border-gray-600">
					<button
						onClick={() => { logoutStaff(); navigate('/login'); }}
						className="flex items-center gap-3 px-3 py-2 rounded-lg text-gray-400 hover:bg-gray-700 hover:text-white w-full font-medium"
					>
						<span className="material-symbols-outlined text-xl">logout</span>
						<span>Cerrar sesión</span>
					</button>
				</div>
			</aside>

			{/* Main Content */}
			<div className="flex-1 flex flex-col">
				{/* Top Header Bar */}
				<header className="flex items-center justify-between p-4 bg-[#1a2831] border-b border-gray-600">
					<div className="relative w-full max-w-md">
						<span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">search</span>
						<input 
							className="w-full pl-10 pr-4 py-2 rounded-lg bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 dark:text-gray-100" 
							placeholder="Buscar componentes..." 
							type="text"
						/>
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
				<main className="flex-1 p-6 space-y-6 bg-[#101c22]">
					{viewMode === 'choose' && (
						<div>
							{/* Page Heading */}
							<div className="flex flex-wrap justify-between gap-3 mb-6">
								<p className="text-gray-900 dark:text-gray-100 text-4xl font-black leading-tight tracking-[-0.033em] min-w-72">¡Bienvenido, Administrador!</p>
							</div>

							{/* Stats */}
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-[#1a2831] border border-gray-600">
									<p className="text-gray-400 text-base font-medium leading-normal">Total Software</p>
									<p className="text-white tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.software || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-[#1a2831] border border-gray-600">
									<p className="text-gray-400 text-base font-medium leading-normal">Total Procesadores</p>
									<p className="text-white tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.procesador || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-[#1a2831] border border-gray-600">
									<p className="text-gray-400 text-base font-medium leading-normal">Total Placas Gráficas</p>
									<p className="text-white tracking-light text-2xl font-bold leading-tight">
										{statsLoading ? '...' : (statsData.placaGrafica || 0)}
									</p>
								</div>
								<div className="flex flex-col gap-2 rounded-lg p-6 bg-[#1a2831] border border-gray-600">
									<p className="text-gray-400 text-base font-medium leading-normal">Total Almacenamiento</p>
									<p className="text-white tracking-light text-2xl font-bold leading-tight">
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

							{/* Categories Grid */}
							<div>
								<div className="flex justify-between items-center pb-3 pt-5">
									<h2 className="text-gray-900 dark:text-gray-100 text-[22px] font-bold leading-tight tracking-[-0.015em]">Categorías de Gestión</h2>
								</div>
								<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
									{entities.map((t) => {
										const count = statsLoading ? '...' : (statsData[t.key] || 0);
										const hasItems = !statsLoading && (statsData[t.key] || 0) > 0;
										
										return (
											<div
												key={t.key}
												className="group relative bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-200 cursor-pointer hover:shadow-lg"
												onClick={() => { setSelectedType(t.key); setForm(emptyTemplate(t.key)); setEditingId(null); }}
											>
												<div className="p-6">
													<div className="flex items-center justify-between mb-4">
														<span className="material-symbols-outlined text-4xl text-blue-600 dark:text-blue-400">{t.icon}</span>
														<div className="text-right">
															<div className="text-lg font-bold text-gray-900 dark:text-gray-100">{count}</div>
															<div className="text-xs text-gray-500 dark:text-gray-400">
																{statsLoading ? 'Cargando...' : hasItems ? 'elementos' : 'Sin elementos'}
															</div>
														</div>
													</div>
													<h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">{t.label}</h3>
													<p className="text-gray-600 dark:text-gray-400 text-sm mb-4">{t.description}</p>
													
													{!statsLoading && !hasItems && (
														<div className="mb-4 p-2 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded text-yellow-700 dark:text-yellow-400 text-xs">
															🚫 Aún no hay elementos cargados
														</div>
													)}
													
													<div className="flex gap-2">
														<button
															className="flex-1 px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors flex items-center justify-center gap-1"
															onClick={(e) => { e.stopPropagation(); setSelectedType(t.key); setForm(emptyTemplate(t.key)); setViewMode('form'); }}
														>
															<span className="material-symbols-outlined text-sm">add</span>
															Agregar
														</button>
														<button
															className={`flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center justify-center gap-1 ${
																hasItems 
																	? 'bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300'
																	: 'bg-gray-50 dark:bg-gray-800 text-gray-400 dark:text-gray-500 cursor-not-allowed'
															}`}
															onClick={(e) => { e.stopPropagation(); if (hasItems) { setSelectedType(t.key); setViewMode('list'); } }}
															disabled={!hasItems}
														>
															<span className="material-symbols-outlined text-sm">list</span>
															Ver Lista
														</button>
													</div>
												</div>
											</div>
										);
									})}
								</div>
							</div>
						</div>
					)}

					{viewMode === 'form' && (
						<div className="max-w-2xl">
							<div className="mb-6">
								<button 
									className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-medium transition-colors"
									onClick={() => setViewMode('choose')}
								>
									<span className="material-symbols-outlined">arrow_back</span>
									Volver
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

									<div className="flex gap-4">
										<button 
											type="submit" 
											disabled={loading}
											className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white rounded-lg font-medium transition-colors"
										>
											{loading ? 'Guardando...' : (editingId ? "Guardar" : "Crear")}
										</button>
										<button
											type="button"
											disabled={loading}
											className="flex-1 px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-medium transition-colors"
											onClick={() => {
												setForm(emptyTemplate(selectedType));
												setEditingId(null);
												setError(null);
											}}
										>
											Cancelar
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
									<p className="text-gray-600 dark:text-gray-400">Total de items: {items.length}</p>
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
												{items.map((it) => (
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