import { ref, computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import type { Project, CreateProjectDto, UpdateProjectDto } from '@/types/project'

export function useProject() {
  const projectStore = useProjectStore()
  const loading = ref(false)

  const currentProject = computed(() => projectStore.currentProject)
  const projects = computed(() => projectStore.projects)

  const loadProjects = async () => {
    loading.value = true
    try {
      await projectStore.loadProjects()
    } finally {
      loading.value = false
    }
  }

  const createProject = async (dto: CreateProjectDto) => {
    loading.value = true
    try {
      return await projectStore.createProject(dto)
    } finally {
      loading.value = false
    }
  }

  const updateProject = async (id: string, dto: UpdateProjectDto) => {
    loading.value = true
    try {
      return await projectStore.updateProject(id, dto)
    } finally {
      loading.value = false
    }
  }

  const deleteProject = async (id: string) => {
    loading.value = true
    try {
      await projectStore.deleteProject(id)
    } finally {
      loading.value = false
    }
  }

  const getProject = async (id: string) => {
    loading.value = true
    try {
      return await projectStore.getProject(id)
    } finally {
      loading.value = false
    }
  }

  const setCurrentProject = (project: Project | null) => {
    projectStore.setCurrentProject(project)
  }

  return {
    loading,
    currentProject,
    projects,
    loadProjects,
    createProject,
    updateProject,
    deleteProject,
    getProject,
    setCurrentProject
  }
}
