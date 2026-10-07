import { DashboardPage } from "../pages/DashboardPage";
import { TurnaroundsPage } from "../pages/TurnaroundsPage";
import { TasksPage } from "../pages/TasksPage";
import { ResourcesPage } from "../pages/ResourcesPage";
import { DelaysPage } from "../pages/DelaysPage";

export const routes = [
  { name: "过站运行看板", route: "/dashboard", page: DashboardPage },
  { name: "航班过站", route: "/turnarounds", page: TurnaroundsPage },
  { name: "地勤任务", route: "/tasks", page: TasksPage },
  { name: "资源调度", route: "/resources", page: ResourcesPage },
  { name: "延误归因", route: "/delays", page: DelaysPage }
] as const;

export type AppRoute = (typeof routes)[number];
